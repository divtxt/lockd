package raftlock

import (
	"github.com/divtxt/lockd/backend"
	"github.com/divtxt/raft"
)

type raftLockImpl struct {
	icmaco  raft.IConsensusModule_AppendCommandOnly
	backend backend.LockBackend
}

// newRaftLockImpl creates a new raftLockImpl with the given params.
//
// The given raft ConsensusModule is assumed to be wired to a BackendAdapter
// that wraps the same LockBackend as the one given here.
//
func newRaftLockImpl(
	icmaco raft.IConsensusModule_AppendCommandOnly,
	backend backend.LockBackend,
) *raftLockImpl {
	return &raftLockImpl{icmaco, backend}
}

func (rli *raftLockImpl) IsLocked(name string) bool {
	return rli.backend.IsLocked(name)
}

func (rli *raftLockImpl) Lock(name string) (bool, error) {
	return rli.lockAction(true, name)
}

func (rli *raftLockImpl) Unlock(name string) (bool, error) {
	return rli.lockAction(false, name)
}

func (rli *raftLockImpl) lockAction(lock bool, name string) (bool, error) {
	// serialize command
	command, err := lockActionSerialize(&lockAction{lock, name})
	if err != nil {
		panic(err) // FIXME: non-panic error handling
	}

	// send to icmaco
	crc, err := rli.icmaco.AppendCommand(command)
	if err != nil {
		return false, err
	}

	// wait for result
	cr, ok := <-crc
	if !ok {
		return false, ErrResultUnknown
	}

	// unpack result
	lar := cr.(lockActionResult)

	return lar.success, nil
}
