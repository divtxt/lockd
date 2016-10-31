package raftlock_test

import (
	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/raft"
	raft_log "github.com/divtxt/raft/log"
	"testing"
)

type Simple_ICM_AppendCommand struct {
	rl     raft.Log
	termNo raft.TermNo
}

func (sicmac *Simple_ICM_AppendCommand) AppendCommand(command raft.Command) (raft.LogIndex, error) {
	entry := raft.LogEntry{sicmac.termNo, command}
	return sicmac.rl.AppendEntry(entry)
}

func TestRaftLock(t *testing.T) {
	var raftLog raft.Log = raft_log.NewInMemoryLog()
	for logIndex := 1; logIndex <= 101; logIndex++ {
		entry := raft.LogEntry{raft.TermNo(logIndex / 10), []byte{}}
		raftLog.AppendEntry(entry)
	}
	iole, err := raftLog.GetIndexOfLastEntry()
	if err != nil {
		t.Fatal(err)
	}
	if iole != 101 {
		t.Fatal(iole)
	}

	raftICMAC := &Simple_ICM_AppendCommand{raftLog, 11}

	rl := raftlock.NewRaftLock(
		raftICMAC,
		raftLog,
		[]string{"bar"},
		101,
	)
	rlCommitApplier := rl.TestHelperGetCommitApplier()
	rlCommitApplier.StopSync()
	rlCommitApplier.TestHelperFakeRestart()

	checkLockState := func(name string, ecs bool, eucs bool) {
		cs, ucs := rl.IsLocked(name)
		if cs != ecs || ucs != eucs {
			t.Fatalf("IsLocked(\"%v\") = (%v, %v) but expected (%v, %v)", name, cs, ucs, ecs, eucs)
		}
	}

	// check starting states
	checkLockState("foo", false, false)
	checkLockState("bar", true, true)

	// unlock "foo" should return nil to indicate unlock failure
	if rl.Unlock("foo") != nil {
		t.Fatal()
	}
	checkLockState("foo", false, false)

	// Lock "foo"
	lockFooCommitChan := rl.Lock("foo")
	if lockFooCommitChan == nil {
		t.Fatal()
	}
	chanWillBlock(t, lockFooCommitChan)
	checkLockState("foo", false, true)

	// a second lock "foo" should return nil to indicate lock failure
	if rl.Lock("foo") != nil {
		t.Fatal()
	}
	checkLockState("foo", false, true)

	// lock "bar" should return nil to indicate lock failure
	if rl.Lock("bar") != nil {
		t.Fatal()
	}
	checkLockState("bar", true, true)

	// Unlock "bar"
	unlockBarCommitChan := rl.Unlock("bar")
	if unlockBarCommitChan == nil {
		t.Fatal()
	}
	chanWillBlock(t, unlockBarCommitChan)
	checkLockState("bar", true, false)

	// Advance commitIndex by 1 log entry
	if rlCommitApplier.TestHelperRunOnceIfTriggerPending() {
		t.Fatal()
	}
	rl.CommitIndexChanged(102)
	if !rlCommitApplier.TestHelperRunOnceIfTriggerPending() {
		t.Fatal()
	}
	checkLockState("foo", true, true)
	checkLockState("bar", true, false) // FAIL
	chanHasValue(t, lockFooCommitChan)
	chanWillBlock(t, unlockBarCommitChan)

	// Relock "bar"
	relockBarCommitChan := rl.Lock("bar")
	if relockBarCommitChan == nil {
		t.Fatal()
	}
	checkLockState("bar", true, true)
	chanWillBlock(t, relockBarCommitChan)

	// Advance commitIndex by 2 log entries
	rl.CommitIndexChanged(104)
	if !rlCommitApplier.TestHelperRunOnceIfTriggerPending() {
		t.Fatal()
	}
	checkLockState("foo", true, true)
	checkLockState("bar", true, true)
	chanHasValue(t, unlockBarCommitChan)
	chanHasValue(t, relockBarCommitChan)
}

func chanWillBlock(t *testing.T, c <-chan struct{}) {
	select {
	case <-c:
		t.Fatal()
	default:
	}
}

func chanHasValue(t *testing.T, c <-chan struct{}) {
	select {
	case <-c:
	default:
		t.Fatal()
	}
}
