package httpimpl

type LockApi interface {
	IsLocked(name string) (bool, bool)
	Lock(name string) (<-chan struct{}, error)
	Unlock(name string) (<-chan struct{}, error)
}
