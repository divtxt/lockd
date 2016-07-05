package locking

import (
	"github.com/divtxt/raft"
	"log"
)

// In-memory implementation of Lock API
type InMemoryLock struct {
	// Lock state
	committedLocks   map[string]bool
	uncommittedLocks map[string]bool

	// Raft things
	raftCM  raft.IConsensusModule
	raftLog raft.LogReadOnly

	// entries                  []raft.LogEntry
	// commitIndex              raft.LogIndex
	// maxEntriesPerAppendEntry uint64
}

// Construct a new InMemoryLock.
//
// Rebuilds the state machine by replaying from the log. This could take a while!
//
// It is expected that raft.ConsensusModule.Start() will be called later using the
// returned InMemoryLock as the raft.ChangeListener parameter.
//
func NewInMemoryLock(raftCM raft.IConsensusModule, raftLog raft.LogReadOnly) *InMemoryLock {
	iml := &InMemoryLock{
		make(map[string]bool),
		make(map[string]bool),
		raftCM,
		raftLog,
	}

	// Build state machine from log
	log.Print("NewInMemoryLock(): Rebuilding state machine from raft log...")
	// iole, err := raftLog.GetIndexOfLastEntry()

	// FIXME: ...

	log.Printf("NewInMemoryLock(): Rebuilt state machine by replaying %v log entries", -1)

	return iml
}

// -- Implement LockApi interface

func (iml *InMemoryLock) Lock(name string) (bool, error) {
	// FIXME: mutex

	// check uncommitted state - if already locked return false
	if _, hasKey := iml.uncommittedLocks[name]; hasKey {
		return false, nil
	}

	// add lock to uncommitted
	iml.uncommittedLocks[name] = true

	// append command to raft log
	command, err := CmdSerialize(&Cmd{true, name})
	if err != nil {
		return false, err
	}
	err = iml.raftCM.AppendCommand(command)
	if err != nil {
		return false, err
	}

	// FIXME: wait for commit

	return true, nil
}

func (iml *InMemoryLock) Unlock(name string) (bool, error) {
	// FIXME: mutex

	// check uncommitted state - if not locked return false
	if _, hasKey := iml.uncommittedLocks[name]; !hasKey {
		return false, nil
	}

	// remove lock from uncommitted
	delete(iml.uncommittedLocks, name)

	// append command to raft log
	command, err := CmdSerialize(&Cmd{false, name})
	if err != nil {
		return false, err
	}
	err = iml.raftCM.AppendCommand(command)
	if err != nil {
		return false, err
	}

	// FIXME: wait for commit

	return true, nil
}

// -- Implement raft.ChangeListener interface

func (iml *InMemoryLock) CommitIndexChanged(logIndex raft.LogIndex) {
	// FIXME: advance commit!
}
