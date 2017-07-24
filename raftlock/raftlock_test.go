package raftlock_test

import (
	"testing"

	"fmt"

	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/raft"
)

func TestRaftLock(t *testing.T) {
	// mock log
	ml := &mockROLog{}
	ml.expectIoleCall(103)
	ml.expectEaiCall(101, []raft.LogEntry{
		raft.LogEntry{5, makeLockCmd("woof")}, // 102
		raft.LogEntry{5, makeLockCmd("moo")},  // 103
		// future entries:
		// 104 = lock foo
		// 105 = unlock bar
		// 106 = unlock woof
		// 107 = unlock moo
		// (discard 106 & 107)
		// 106 = relock bar
		// 107 = unlock woof
		// 108 = lock tweet
	})

	// Create RaftLock instance.
	rl := raftlock.NewRaftLock(
		ml,
		[]string{"bar"},
		101)
	if rl.GetLastApplied() != 101 {
		t.Fatal(rl.GetLastApplied())
	}
	if !ml.ioleCalled || !ml.eaiCalled {
		t.Fatal(ml)
	}

	//
	checkLockState := func(name string, ecs bool, eucs bool) {
		cs, ucs := rl.IsLocked(name)
		if cs != ecs || ucs != eucs {
			panic(fmt.Sprintf(
				"IsLocked(\"%v\") = (%v, %v) but expected (%v, %v)", name, cs, ucs, ecs, eucs,
			))
		}
	}

	// check starting states, and verify newer entries from Log were applied
	checkLockState("foo", false, false)
	checkLockState("bar", true, true)
	checkLockState("woof", false, true)
	checkLockState("moo", false, true)

	// unlock "foo" should return ErrAlreadyUnlocked to indicate unlock failure
	err := rl.CheckAndApplyCommand(0, makeUnlockCmd("foo"))
	if err != raftlock.ErrAlreadyUnlocked {
		t.Fatal(err)
	}
	checkLockState("foo", false, false)

	// -- 104 : Lock "foo"
	err = rl.CheckAndApplyCommand(104, makeLockCmd("foo"))
	if err != nil {
		t.Fatal(err)
	}
	checkLockState("foo", false, true)

	// a second lock "foo" should return ErrAlreadyLocked to indicate lock failure
	err = rl.CheckAndApplyCommand(0, makeLockCmd("foo"))
	if err != raftlock.ErrAlreadyLocked {
		t.Fatal(err)
	}
	checkLockState("foo", false, true)

	// lock "bar" should return nil to indicate lock failure
	err = rl.CheckAndApplyCommand(0, makeLockCmd("bar"))
	if err != raftlock.ErrAlreadyLocked {
		t.Fatal(err)
	}
	checkLockState("bar", true, true)

	// -- 105 : Unlock "bar"
	err = rl.CheckAndApplyCommand(105, makeUnlockCmd("bar"))
	if err != nil {
		t.Fatal(err)
	}
	checkLockState("bar", true, false)

	// -- Commit to 104 : Lock "foo"
	ml.expectIoleCall(104)
	ml.expectEaiCall(101, []raft.LogEntry{
		raft.LogEntry{5, makeLockCmd("woof")}, // 102
		raft.LogEntry{5, makeLockCmd("moo")},  // 103
		raft.LogEntry{5, makeLockCmd("foo")},  // 104
	})
	rl.CommitIndexChanged(104)
	if ml.ioleCalled || !ml.eaiCalled {
		t.Fatal(ml)
	}
	checkLockState("foo", true, true)
	checkLockState("bar", true, false)
	checkLockState("woof", true, true)
	checkLockState("moo", true, true)

	// -- 106 & 107 : Unlock "woof" and "moo"
	err = rl.CheckAndApplyCommand(106, makeUnlockCmd("woof"))
	if err != nil {
		t.Fatal(err)
	}
	err = rl.CheckAndApplyCommand(107, makeUnlockCmd("moo"))
	if err != nil {
		t.Fatal(err)
	}
	checkLockState("foo", true, true)
	checkLockState("bar", true, false)
	checkLockState("woof", true, false)
	checkLockState("moo", true, false)

	// -- redo 106 & 107, add 108 : Discard & replay 3 log entries
	ml.expectIoleCall(108)
	ml.expectEaiCall(104, []raft.LogEntry{
		raft.LogEntry{5, makeUnlockCmd("bar")},  // 105
		raft.LogEntry{5, makeLockCmd("bar")},    // 106
		raft.LogEntry{5, makeUnlockCmd("woof")}, // 107
		raft.LogEntry{5, makeLockCmd("tweet")},  // 108
	})
	err = rl.SetEntriesAfterIndex(105, []raft.LogEntry{
		raft.LogEntry{5, makeLockCmd("bar")},    // 106
		raft.LogEntry{5, makeUnlockCmd("woof")}, // 107
		raft.LogEntry{5, makeLockCmd("tweet")},  // 108
	})
	if err != nil {
		t.Fatal(err)
	}
	if !ml.ioleCalled || !ml.eaiCalled {
		t.Fatal(ml)
	}
	checkLockState("foo", true, true)
	checkLockState("bar", true, true)
	checkLockState("woof", true, false)
	checkLockState("moo", true, true)
	checkLockState("tweet", false, true)

	// -- Commit to 105 : Unlock "bar"
	ml.expectIoleCall(105)
	ml.expectEaiCall(104, []raft.LogEntry{
		raft.LogEntry{5, makeUnlockCmd("bar")}, // 105
	})
	rl.CommitIndexChanged(105)
	if ml.ioleCalled || !ml.eaiCalled {
		t.Fatal(ml)
	}
	checkLockState("foo", true, true)
	checkLockState("bar", false, true)
	checkLockState("woof", true, false)
	checkLockState("moo", true, true)
	checkLockState("tweet", false, true)
}

// -- mock LogReadOnly

func makeLockCmd(n string) raft.Command {
	cmd, err := raftlock.MakeLockCommand(n)
	if err != nil {
		panic(err)
	}
	return cmd
}

func makeUnlockCmd(n string) raft.Command {
	cmd, err := raftlock.MakeUnlockCommand(n)
	if err != nil {
		panic(err)
	}
	return cmd
}

type mockROLog struct {
	ioleCalled bool
	iole       raft.LogIndex

	// expectedTaiIndex raft.LogIndex
	// tai              raft.TermNo

	eaiCalled        bool
	expectedEaiIndex raft.LogIndex
	expectedEaiCount uint64
	eai              []raft.LogEntry
}

func (m *mockROLog) expectIoleCall(i raft.LogIndex) {
	m.ioleCalled = false
	m.iole = i
}

func (m *mockROLog) GetIndexOfLastEntry() (raft.LogIndex, error) {
	m.ioleCalled = true
	return m.iole, nil
}

func (m *mockROLog) GetTermAtIndex(i raft.LogIndex) (raft.TermNo, error) {
	panic("Not implemented!")
}

func (m *mockROLog) expectEaiCall(i raft.LogIndex, e []raft.LogEntry) {
	m.eaiCalled = false
	m.expectedEaiIndex = i
	m.expectedEaiCount = uint64(len(e))
	m.eai = e
}

func (m *mockROLog) GetEntriesAfterIndex(i raft.LogIndex, c uint64) ([]raft.LogEntry, error) {
	if m.eaiCalled {
		panic("double call!")
	}
	if i != m.expectedEaiIndex {
		panic(fmt.Sprintf("%v", i))
	}
	if c != m.expectedEaiCount {
		panic(fmt.Sprintf("%v", c))
	}
	m.eaiCalled = true
	return m.eai, nil
}
