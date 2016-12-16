package statemachine

// Interface to a Lock state machine implementation.
//
// Concurrence safety is not required by this interface.
//
type LockStateMachine interface {
	IsLocked(name string) bool
	Lock(name string) bool
	Unlock(name string) bool
}
