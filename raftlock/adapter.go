package raftlock

import (
	"github.com/divtxt/lockd/backend"
	"github.com/divtxt/raft"
)

// smAdapter wraps a LockBackend to implement the raft.StateMachine interface.
type smAdapter struct {
	backend backend.LockBackend
}

func (sma smAdapter) GetLastApplied() raft.LogIndex {
	return sma.backend.GetLastApplied()
}

func (sma smAdapter) ApplyCommand(
	logIndex raft.LogIndex,
	command raft.Command,
) raft.CommandResult {
	// Deserialize the command
	cmd, err := lockActionDeserialize(command)
	if err != nil {
		panic(err)
	}

	// Apply to the backend
	var success bool
	if cmd.Lock {
		success = sma.backend.Lock(logIndex, cmd.Name)
	} else {
		success = sma.backend.Unlock(logIndex, cmd.Name)
	}

	return lockActionResult{success}
}
