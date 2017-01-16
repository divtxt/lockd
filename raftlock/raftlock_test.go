package raftlock_test

import (
	"testing"

	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/raft"
	raft_committer "github.com/divtxt/raft/committer"
	raft_log "github.com/divtxt/raft/log"
)

func TestRaftLock(t *testing.T) {
	// Real log for testing
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

	// Create RaftLock instance.
	rl := raftlock.NewRaftLock(
		raftLog,
		[]string{"bar"},
		101,
	)

	// Create (partial) ConsensusModule.
	// Simplify testing by not using a full ConsensusModule but by using some parts.
	micmaco := NewMockICMACO(raftLog, 11, rl)
	// Use a real Committer in test mode to drive commits.
	committer := raft_committer.NewCommitter(raftLog, rl)
	committer.StopSync() // switch to manual control
	committer.TestHelperGetCommitApplier().TestHelperFakeRestart()

	// Give RaftLock the IConsensusModule_AppendCommandOnly reference.
	rl.SetICMACO(micmaco)

	//
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
	if committer.TestHelperGetCommitApplier().TestHelperRunOnceIfTriggerPending() {
		t.Fatal()
	}
	committer.CommitIndexChanged(102)
	if !committer.TestHelperGetCommitApplier().TestHelperRunOnceIfTriggerPending() {
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
	committer.CommitIndexChanged(104)
	if !committer.TestHelperGetCommitApplier().TestHelperRunOnceIfTriggerPending() {
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

// ---- Mock IConsensusModule_AppendCommandOnly

type MockICMACO struct {
	rl     raft.Log
	termNo raft.TermNo
	sm     raft.StateMachine
}

func NewMockICMACO(rl raft.Log, termNo raft.TermNo, sm raft.StateMachine) *MockICMACO {
	return &MockICMACO{rl, termNo, sm}
}

func (micmaco *MockICMACO) AppendCommand(command raft.Command) (raft.LogIndex, error) {
	entry := raft.LogEntry{micmaco.termNo, command}
	return micmaco.rl.AppendEntry(entry)
}
