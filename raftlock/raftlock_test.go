package raftlock_test

import (
	"fmt"
	"testing"

	"github.com/divtxt/lockd/backend"
	"github.com/divtxt/lockd/raftlock"
	"github.com/divtxt/raft"
)

func TestRaftLock(t *testing.T) {
	imb := backend.NewInMemoryBackend(10)
	sm, prl := raftlock.NewRaftLock(imb)
	micmaco := newMockICMACO(sm)
	rl := prl.Finish(micmaco)

	rl_IsLocked := func(name string, expected bool) {
		if actual := rl.IsLocked(name); actual != expected {
			panic(fmt.Sprintf("%v", actual))
		}
	}

	rl_Lock := func(name string, expected bool, expectedErr error) {
		success, err := rl.Lock(name)
		if err != expectedErr {
			panic(err)
		}
		if expectedErr == nil {
			if success != expected {
				panic(fmt.Sprintf("got %v / expected %v", success, expected))
			}
		}
	}

	rl_Unlock := func(name string, expected bool, expectedErr error) {
		success, err := rl.Unlock(name)
		if err != expectedErr {
			panic(err)
		}
		if expectedErr == nil {
			if success != expected {
				panic(fmt.Sprintf("got %v / expected %v", success, expected))
			}
		}
	}
	// Initial state test
	rl_IsLocked("foo", false)
	rl_IsLocked("bar", false)

	// Check that errors are returned
	micmaco.sendErr = raft.ErrStopped
	rl_Lock("foo", false, raft.ErrStopped)
	micmaco.sendErr = nil

	// Lock some entries
	rl_Lock("foo", true, nil)
	rl_Lock("bar", true, nil)
	rl_IsLocked("foo", true)
	rl_IsLocked("bar", true)

	// Lock of locked entry should fail
	rl_Lock("foo", false, nil)
	rl_IsLocked("foo", true)

	// Unlock one entry
	rl_Unlock("foo", true, nil)
	rl_IsLocked("foo", false)
	rl_IsLocked("bar", true)

	// Unlock again should fail
	rl_Unlock("foo", false, nil)
	rl_IsLocked("foo", false)

	// Unlock another entry, lose leadership and also lose the action
	micmaco.loseLeadershipAndAction = true
	rl_Unlock("bar", false, raftlock.ErrResultUnknown)
	micmaco.loseLeadershipAndAction = false
	rl_IsLocked("bar", true)

	// Unlock again, lose leadership but commit the action
	micmaco.loseLeadershipButCommitAction = true
	rl_Unlock("bar", false, raftlock.ErrResultUnknown)
	micmaco.loseLeadershipButCommitAction = false
	rl_IsLocked("bar", false)
}

// ---- Mock ConsensusModule - only need to implement IConsensusModule_AppendCommandOnly

type mockICMACO struct {
	sendErr                       error
	loseLeadershipAndAction       bool
	loseLeadershipButCommitAction bool
	commitIndex                   raft.LogIndex
	sm                            raft.StateMachine
}

func newMockICMACO(sm raft.StateMachine) *mockICMACO {
	return &mockICMACO{nil, false, false, sm.GetLastApplied(), sm}
}

func (micmaco *mockICMACO) AppendCommand(command raft.Command) (<-chan raft.CommandResult, error) {
	if micmaco.sendErr != nil {
		return nil, micmaco.sendErr
	}
	micmaco.commitIndex++
	crc := make(chan raft.CommandResult, 1)
	if micmaco.loseLeadershipAndAction {
		close(crc)
		return crc, nil
	}
	cr := micmaco.sm.ApplyCommand(micmaco.commitIndex, command)
	if micmaco.loseLeadershipButCommitAction {
		close(crc)
		return crc, nil
	}
	crc <- cr
	return crc, nil
}
