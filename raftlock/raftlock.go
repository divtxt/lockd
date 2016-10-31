package raftlock

import (
	"fmt"
	"github.com/divtxt/lockd/runner"
	"github.com/divtxt/lockd/statemachine"
	"github.com/divtxt/raft"
)

// Raft-commit driven lock.
type RaftLock struct {
	// Raft things
	raftICMAC ICM_AppendCommand
	raftLog   raft.LogReadOnly

	// Lock state
	committedLocks   statemachine.LockStateMachine
	uncommittedLocks statemachine.LockStateMachine

	// Commit state
	commitIndex        raft.LogIndex
	appliedCommitIndex raft.LogIndex
	commitApplier      *runner.TriggeredRunner

	// Reply channels
	replyChans map[raft.LogIndex]chan struct{}
}

// The subset of the Consensus interface that RaftLock cares about:
type ICM_AppendCommand interface {
	AppendCommand(command raft.Command) (raft.LogIndex, error)
}

// Construct a new RaftLock.
//
// It is expected that raft.ConsensusModule.Start() will be called later using the
// returned RaftLock as the raft.ChangeListener parameter.
//
func NewRaftLock(
	raftICMAC ICM_AppendCommand,
	raftLog raft.LogReadOnly,
	initialLocks []string,
	initialCommitIndex raft.LogIndex,
) *RaftLock {

	// TODO: copy initial state
	rl := &RaftLock{
		raftICMAC,
		raftLog,
		statemachine.NewInMemoryLSM(),
		statemachine.NewInMemoryLSM(),
		initialCommitIndex,
		initialCommitIndex,
		nil, // commitApplier
		make(map[raft.LogIndex]chan struct{}),
	}

	err := applyLocks(initialLocks, rl.committedLocks)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	err = applyLocks(initialLocks, rl.uncommittedLocks)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	// Start commitApplier goroutine
	rl.commitApplier = runner.NewTriggeredRunner(rl.applyPendingCommits)

	return rl
}

// ---- API

// Check if the given entry is locked.
//
// This returns a pair of booleans indicating the committed and uncommitted states of the entry.
//
func (rl *RaftLock) IsLocked(name string) (bool, bool) {
	// FIXME: mutex

	// check uncommitted state - if already locked return nil
	committedState, err := rl.committedLocks.IsLocked(name)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	uncommittedState, err := rl.uncommittedLocks.IsLocked(name)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	return committedState, uncommittedState
}

// Lock the given entry.
//
// If the entry is unlocked, a lock action for the entry is appended to the raft log,
// and a channel is returned.
// The call does not wait for the entry in the raft log to successfully commit
// to the raft cluster.
// A dummy value is sent later on the returned channel when the entry is committed.
// If the entry does not commit, no value is ever sent on the returned channel.
//
// If the entry is already locked, a nil is returned.
//
// Note that uncommitted lock state is tracked to ensure correct behavior.
// For example: a locked entry can be unlocked even if the lock action has not yet committed.
// Or: an entry cannot be locked twice even if the first lock action has not yet committed.
//
func (rl *RaftLock) Lock(name string) <-chan struct{} {
	// FIXME: mutex

	// check uncommitted state - if already locked return nil
	alreadyLocked, err := rl.uncommittedLocks.IsLocked(name)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	if alreadyLocked {
		return nil
	}

	// append command to raft log
	command, err := CmdSerialize(&Cmd{true, name})
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	logIndex, err := rl.raftICMAC.AppendCommand(command)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	// apply lock action to uncommitted
	lockSuccess, err := rl.uncommittedLocks.Lock(name)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	if !lockSuccess {
		panic(fmt.Sprintf("FATAL: uncommittedLocks.Lock() unexpectedly already locked: %v", name))
	}

	return rl.makeReplyChan(logIndex)
}

// Unlock the given entry.
//
// If the entry is locked, an unlock action for the entry is appended to the raft log,
// and a channel is returned.
// The call does not wait for the entry in the raft log to successfully commit
// to the raft cluster.
// A dummy value is sent later on the returned channel when the entry is committed.
// If the entry does not commit, no value is ever sent on the returned channel.
//
// If the entry is not locked, a nil is returned.
//
// Note that uncommitted lock state is tracked to ensure correct behavior.
// For example: an unlocked entry can be locked even if the unlock action has not yet committed.
// Or: an entry cannot be unlocked twice even if the first unlock action has not yet committed.
//
func (rl *RaftLock) Unlock(name string) <-chan struct{} {
	// FIXME: mutex

	// check uncommitted state - if not locked return nil
	alreadyLocked, err := rl.uncommittedLocks.IsLocked(name)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	if !alreadyLocked {
		return nil
	}

	// append command to raft log
	command, err := CmdSerialize(&Cmd{false, name})
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	logIndex, err := rl.raftICMAC.AppendCommand(command)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	// apply unlock action to uncommitted
	unlockSuccess, err := rl.uncommittedLocks.Unlock(name)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	if !unlockSuccess {
		panic(fmt.Sprintf("FATAL: uncommittedLocks.Unlock() unexpectedly not locked: %v", name))
	}

	return rl.makeReplyChan(logIndex)
}

// Receive raft commit index changes.
//
// Commit are applied asynchronously.
func (rl *RaftLock) CommitIndexChanged(commitIndex raft.LogIndex) {
	// FIXME: mutex

	// Check commitIndex is not going backward
	if rl.commitIndex > commitIndex {
		panic(fmt.Sprintf("FATAL: decreasing commitIndex: %v > %v", rl.commitIndex, commitIndex))
	}

	// Trigger application of committed entries
	rl.commitApplier.TriggerRun()

	// Update commitIndex
	rl.commitIndex = commitIndex
}

// ---- Private

// Lock the listed entries.
//
// Failures to lock for entries already locked or duplicate are ignored.
func applyLocks(locks []string, lsm statemachine.LockStateMachine) error {
	for _, l := range locks {
		_, err := lsm.Lock(l)
		if err != nil {
			return err
		}
	}
	return nil
}

//
func (rl *RaftLock) makeReplyChan(logIndex raft.LogIndex) <-chan struct{} {
	replyChan := make(chan struct{}, 1)
	if _, hasKey := rl.replyChans[logIndex]; hasKey {
		panic(fmt.Sprintf("FATAL: replyChan already exists for logIndex=%v", logIndex))
	}
	rl.replyChans[logIndex] = replyChan
	return replyChan
}

//
func (rl *RaftLock) sendReply(logIndex raft.LogIndex) {
	replyChan, hasKey := rl.replyChans[logIndex]
	if !hasKey {
		panic(fmt.Sprintf("FATAL: no replyChan for logIndex=%v", logIndex))
	}
	delete(rl.replyChans, logIndex)
	replyChan <- struct{}{}
}

// Apply all pending committed entries and notify callers via replyChans.
//
// Does nothing if there are no committed entries to apply.
func (rl *RaftLock) applyPendingCommits() {
	// FIXME: mutex

	for rl.appliedCommitIndex < rl.commitIndex {
		rl.applyOnePendingCommit()
	}

}

// Apply one commit
func (rl *RaftLock) applyOnePendingCommit() {
	// FIXME: mutex

	if rl.appliedCommitIndex >= rl.commitIndex {
		return
	}

	indexToApply := rl.appliedCommitIndex + 1

	// Get one command from the raft log
	// TODO: get and apply multiple entries at a time
	entries, err := rl.raftLog.GetEntriesAfterIndex(indexToApply-1, 1)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	// Deserialize the command
	cmd, err := CmdDeserialize(entries[0].Command)
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	// Apply to committedLocks
	name := cmd.Name
	if cmd.Lock {
		lockSuccess, err := rl.committedLocks.Lock(name)
		if err != nil {
			panic(err) // FIXME: non-panic error handling
		}
		if !lockSuccess {
			panic(fmt.Sprintf("FATAL: committedLocks.Lock() unexpectedly already locked: %v", name))
		}
	} else {
		unlockSuccess, err := rl.committedLocks.Unlock(name)
		if err != nil {
			panic(err) // FIXME: non-panic error handling
		}
		if !unlockSuccess {
			panic(fmt.Sprintf("FATAL: committedLocks.Unlock() unexpectedly not locked: %v", name))
		}
	}

	// Send reply for rl.appliedCommitIndex
	rl.sendReply(indexToApply)

	// Update state
	rl.appliedCommitIndex = indexToApply
}

// Meant only for tests!
func (rl *RaftLock) TestHelperGetCommitApplier() *runner.TriggeredRunner {
	return rl.commitApplier
}
