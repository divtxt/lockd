package raftlock

import (
	"errors"
	"sync"

	"fmt"

	"github.com/divtxt/lockd/statemachine"
	"github.com/divtxt/raft"
)

var ErrAlreadyLocked = errors.New("Already locked")
var ErrAlreadyUnlocked = errors.New("Already locked")

// RaftLock is a locking service implementing raft.StateMachine interface.
type RaftLock struct {
	mutex *sync.Mutex

	raftLog raft.LogReadOnly

	// Lock state
	committedLocks   *statemachine.InMemoryLSM
	uncommittedLocks *statemachine.InMemoryLSM
	lastApplied      raft.LogIndex
}

// Construct a new RaftLock.
//
func NewRaftLock(
	raftLog raft.LogReadOnly,
	initialLocks []string,
	initialLastApplied raft.LogIndex,
) *RaftLock {

	// TODO: copy initial state
	rl := &RaftLock{
		&sync.Mutex{},
		raftLog,
		statemachine.NewInMemoryLSM(),
		nil,
		initialLastApplied,
	}

	// FIXME: move this to NewInMemoryLSM
	for _, l := range initialLocks {
		rl.committedLocks.Lock(l)
	}

	rl.rebuildUncommittedLocks()

	return rl
}

// ---- API

// IsLocked checks if the given entry is locked.
//
// This returns a pair of booleans indicating the committed and uncommitted states of the entry.
//
func (rl *RaftLock) IsLocked(name string) (bool, bool) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	committedState := rl.committedLocks.IsLocked(name)
	uncommittedState := rl.uncommittedLocks.IsLocked(name)

	return committedState, uncommittedState
}

// ---- Implement raft.StateMachine interface

// GetLastApplied should return the value of lastApplied.
func (rl *RaftLock) GetLastApplied() raft.LogIndex {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	return rl.lastApplied
}

// CheckAndApplyCommand should check if the given command against the state machine
// and either apply it if allowed or return an error if not allowed.
func (rl *RaftLock) CheckAndApplyCommand(
	logIndex raft.LogIndex, command raft.Command,
) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Deserialize the command
	cmd, err := lockActionDeserialize(command)
	if err != nil {
		panic(err)
	}

	// Apply to uncommittedLocks
	name := cmd.Name
	if cmd.Lock {
		lockSuccess := rl.uncommittedLocks.Lock(name)
		if !lockSuccess {
			return ErrAlreadyLocked
		}
	} else {
		unlockSuccess := rl.uncommittedLocks.Unlock(name)
		if !unlockSuccess {
			return ErrAlreadyUnlocked
		}
	}

	return nil
}

// SetEntriesAfterIndex tells the state machine about Log entry changes.
func (rl *RaftLock) SetEntriesAfterIndex(afterIndex raft.LogIndex, entries []raft.LogEntry) error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	iole, err := rl.raftLog.GetIndexOfLastEntry()
	if err != nil {
		// FIXME
		panic(err)
	}

	// Avoid rebuilding if nothing is actually being discarded
	if afterIndex == iole && len(entries) == 0 {
		return nil
	}

	// Rebuild uncommittedLocks
	rl.rebuildUncommittedLocks()

	return nil
}

// CommitIndexChanged tells the state machine that commitIndex has changed to the given value.
func (rl *RaftLock) CommitIndexChanged(newCommitIndex raft.LogIndex) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Apply committed entries to committedLocks
	lastApplied := rl.lastApplied
	if newCommitIndex < lastApplied {
		panic(fmt.Sprintf("newCommitIndex=%v is < lastApplied=%v", newCommitIndex, lastApplied))
	}
	numEntries := uint64(newCommitIndex - lastApplied)
	entries, err := rl.raftLog.GetEntriesAfterIndex(lastApplied, numEntries)
	if err != nil {
		panic(err)
	}
	applyLogEntries(rl.committedLocks, entries)

	// Update lastApplied
	rl.lastApplied = newCommitIndex
}

// ---- Private

func (rl *RaftLock) rebuildUncommittedLocks() {
	rl.uncommittedLocks = rl.committedLocks.Clone()

	iole, err := rl.raftLog.GetIndexOfLastEntry()
	if err != nil {
		panic(err)
	}
	lastApplied := rl.lastApplied

	if iole < lastApplied {
		panic(fmt.Sprintf("indexOfLastEntry=%v is < lastApplied=%v", iole, lastApplied))
	}

	numEntries := uint64(iole - lastApplied)
	entries, err := rl.raftLog.GetEntriesAfterIndex(lastApplied, numEntries)
	if err != nil {
		panic(err)
	}

	applyLogEntries(rl.uncommittedLocks, entries)
}

func applyLogEntries(lsm statemachine.LockStateMachine, entries []raft.LogEntry) {
	for _, e := range entries {
		// Deserialize the command
		cmd, err := lockActionDeserialize(e.Command)
		if err != nil {
			panic(err)
		}
		// Apply it
		name := cmd.Name
		if cmd.Lock {
			lockSuccess := lsm.Lock(name)
			if !lockSuccess {
				panic(fmt.Sprintf("FATAL: unexpectedly already locked: %v", name))
			}
		} else {
			unlockSuccess := lsm.Unlock(name)
			if !unlockSuccess {
				panic(fmt.Sprintf("FATAL: unexpectedly already locked: %v", name))
			}
		}
	}
}
