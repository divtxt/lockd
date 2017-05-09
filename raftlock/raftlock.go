package raftlock

import (
	"fmt"
	"sync"

	"github.com/divtxt/lockd/statemachine"
	"github.com/divtxt/raft"
)

// Raft-commit driven lock.
type RaftLock struct {
	mutex *sync.Mutex

	// Raft things
	raftICMACO raft.IConsensusModule_AppendCommandOnly
	raftLog    raft.LogReadOnly

	// Lock state
	committedLocks   statemachine.LockStateMachine
	uncommittedLocks statemachine.LockStateMachine
	lastApplied      raft.LogIndex

	// Reply channels
	replyChans map[raft.LogIndex]chan struct{}
}

// Construct a new RaftLock.
//
// The RaftLock setup is incomplete until SetICMACO() is called with the raft ConsensusModule.
//
func NewRaftLock(
	raftLog raft.LogReadOnly,
	initialLocks []string,
	initialLastApplied raft.LogIndex,
) *RaftLock {

	// TODO: copy initial state
	rl := &RaftLock{
		&sync.Mutex{},
		nil, // raftICMACO
		raftLog,
		statemachine.NewInMemoryLSM(),
		statemachine.NewInMemoryLSM(),
		initialLastApplied,
		make(map[raft.LogIndex]chan struct{}),
	}

	applyLocks(initialLocks, rl.committedLocks)
	applyLocks(initialLocks, rl.uncommittedLocks)

	return rl
}

func (rl *RaftLock) SetICMACO(raftICMACO raft.IConsensusModule_AppendCommandOnly) {
	if rl.raftICMACO != nil {
		panic("Attempt to set raftICMACO more than once!")
	}
	if raftICMACO == nil {
		panic("SetICMACO() called with nil!")
	}
	rl.raftICMACO = raftICMACO
}

// ---- API

// Check if the given entry is locked.
//
// This returns a pair of booleans indicating the committed and uncommitted states of the entry.
//
func (rl *RaftLock) IsLocked(name string) (bool, bool) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// check uncommitted state - if already locked return nil
	committedState := rl.committedLocks.IsLocked(name)
	uncommittedState := rl.uncommittedLocks.IsLocked(name)

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
func (rl *RaftLock) Lock(name string) (<-chan struct{}, error) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// check uncommitted state - if already locked return nil
	alreadyLocked := rl.uncommittedLocks.IsLocked(name)
	if alreadyLocked {
		return nil, nil
	}

	// append command to raft log
	command, err := lockActionSerialize(&lockAction{true, name})
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	logIndex, err := rl.raftICMACO.AppendCommand(command)
	if err != nil {
		return nil, err
	}

	// apply lock action to uncommitted
	lockSuccess := rl.uncommittedLocks.Lock(name)
	if !lockSuccess {
		panic(fmt.Sprintf("FATAL: uncommittedLocks.Lock() unexpectedly already locked: %v", name))
	}

	return rl.makeReplyChan(logIndex), nil
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
func (rl *RaftLock) Unlock(name string) (<-chan struct{}, error) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// check uncommitted state - if not locked return nil
	alreadyLocked := rl.uncommittedLocks.IsLocked(name)
	if !alreadyLocked {
		return nil, nil
	}

	// append command to raft log
	command, err := lockActionSerialize(&lockAction{false, name})
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}
	logIndex, err := rl.raftICMACO.AppendCommand(command)
	if err != nil {
		return nil, err
	}

	// apply unlock action to uncommitted
	unlockSuccess := rl.uncommittedLocks.Unlock(name)
	if !unlockSuccess {
		panic(fmt.Sprintf("FATAL: uncommittedLocks.Unlock() unexpectedly not locked: %v", name))
	}

	return rl.makeReplyChan(logIndex), nil
}

// ---- Implement raft.StateMachine interface

// GetLastApplied should return the value of lastApplied.
func (rl *RaftLock) GetLastApplied() raft.LogIndex {
	return rl.lastApplied
}

// ApplyCommand should apply the given command to the state machine.
func (rl *RaftLock) ApplyCommand(logIndex raft.LogIndex, command raft.Command) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Check lastApplied is not going backward
	if rl.lastApplied > logIndex {
		panic(fmt.Sprintf("FATAL: decreasing lastApplied: %v > %v", rl.lastApplied, logIndex))
	}

	// Deserialize the command
	cmd, err := lockActionDeserialize(command)
	if err != nil {
		panic(err)
	}

	// Apply to committedLocks
	name := cmd.Name
	if cmd.Lock {
		lockSuccess := rl.committedLocks.Lock(name)
		if !lockSuccess {
			panic(fmt.Sprintf("FATAL: committedLocks.Lock() unexpectedly already locked: %v", name))
		}
	} else {
		unlockSuccess := rl.committedLocks.Unlock(name)
		if !unlockSuccess {
			panic(fmt.Sprintf("FATAL: committedLocks.Unlock() unexpectedly not locked: %v", name))
		}
	}

	// Update lastApplied
	rl.lastApplied = logIndex

	// Send reply for rl.appliedCommitIndex
	rl.sendReply(logIndex)
}

// ---- Private

// Lock the listed entries.
//
// Failures to lock for entries already locked or duplicate are ignored.
func applyLocks(locks []string, lsm statemachine.LockStateMachine) {
	for _, l := range locks {
		lsm.Lock(l)
	}
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
	replyChan := rl.replyChans[logIndex]
	if replyChan != nil {
		delete(rl.replyChans, logIndex)
		replyChan <- struct{}{}
	}
}
