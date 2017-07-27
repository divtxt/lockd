package raftlock

import (
	"errors"
)

// RaftLock is a locking service implemented on raft.
//
// Errors returned by Lock() and Unlock() are one of:
// - raft.ErrStopped if ConsensusModule is stopped.
// - raft.ErrNotLeader if not currently the leader.
// - ErrResultUnknown if leadership was lost during the operation.
//
type RaftLock interface {
	IsLocked(name string) bool
	Lock(name string) (bool, error)
	Unlock(name string) (bool, error)
}

var ErrResultUnknown = errors.New("Result Unknown (leadership was lost before commit)")
