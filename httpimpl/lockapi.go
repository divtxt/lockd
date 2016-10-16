package httpimpl

type LockApi interface {
	Lock(name string) (bool, error)
	Unlock(name string) (bool, error)
}
