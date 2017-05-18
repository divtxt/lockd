package httpimpl

// LockApi is internal interface to raft aware lock state machine.
//
// Errors returned by Lock() and Unlock() are those returned by raft.IConsensusModule.AppendCommand()
// i.e. one of:
// - raft.ErrStopped if ConsensusModule is stopped.
// - raft.ErrNotLeader if not currently the leader.
type LockApi interface {
	IsLocked(name string) (bool, bool)
	Lock(name string) (<-chan struct{}, error)
	Unlock(name string) (<-chan struct{}, error)
}
