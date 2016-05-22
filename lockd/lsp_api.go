package lockd

// Locking internal API to abstract the state persistence implementation.
type LockStatePersistence interface {
	Lock(name string) (bool, error)
	Unlock(name string) (bool, error)
}
