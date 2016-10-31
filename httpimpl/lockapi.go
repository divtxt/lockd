package httpimpl

type LockApi interface {
	IsLocked(name string) (bool, bool)
	Lock(name string) <-chan struct{}
	Unlock(name string) <-chan struct{}
}
