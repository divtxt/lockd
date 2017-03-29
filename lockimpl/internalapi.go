package lockimpl

type InternalLockApi interface {
	IsLocked(name string) (bool, bool)
	Lock(name string) <-chan struct{}
	Unlock(name string) <-chan struct{}
}
