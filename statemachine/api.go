package statemachine

// Interface to a Lock state machine implementation.
//
// Concurrence safety is not required by this interface.
//
type LockStateMachine interface {
	IsLocked(name string) (bool, error)
	Lock(name string) (bool, error)
	Unlock(name string) (bool, error)
}
