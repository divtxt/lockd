package statemachine

// Interface to a Lock state machine implementation.
//
// Implementation notes:
// - Concurrence safety is not required by this interface.
// - Can panic for bad lock names - callers should already have validated names with IsValidName().
//
type LockStateMachine interface {
	IsLocked(name string) bool
	Lock(name string) bool
	Unlock(name string) bool
}
