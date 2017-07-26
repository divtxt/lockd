package raftlock

import (
	"github.com/divtxt/lockd/backend"
	"github.com/divtxt/raft"
)

type PartialRaftLock interface {
	Finish(raft.IConsensusModule_AppendCommandOnly) RaftLock
}

// NewRaftLock starts the creation of a RaftLock with the given LockBackend.
//
// A raft ConsensusModule should be built with the returned StateMachine.
// The ConsensusModule should then be sent to the returned PartialRaftLock to finish construction.
//
func NewRaftLock(backend backend.LockBackend) (raft.StateMachine, PartialRaftLock) {
	ba := smAdapter{backend}
	p := &partial{backend}
	return ba, p
}

type partial struct {
	backend backend.LockBackend
}

func (p *partial) Finish(icmaco raft.IConsensusModule_AppendCommandOnly) RaftLock {
	if p.backend == nil {
		panic("Called more than once!")
	}
	rli := newRaftLockImpl(icmaco, p.backend)
	p.backend = nil
	return rli
}
