package backend

import "github.com/divtxt/raft"

// LockBackend is a lock state machine backend that sits behind raft.
//
// Implementation notes:
// - Behavior should be deterministic for a given current state and given action.
// - IsLocked may be called concurrently with other methods.
// - All other methods are only called one at a time.
// - Can panic for bad lock names - callers should already have validated names with IsValidName().
// - All actions should save the given log index as the new lastApplied value. However, since
// the implementation should be deterministic, actions that do not change the state can choose
// to not save the value. For example, a Lock() call that does not acquire the lock.
//
type LockBackend interface {
	// Get the current value of lastApplied.
	GetLastApplied() raft.LogIndex

	// Check if the given entry is locked.
	IsLocked(name string) bool

	// Lock the given entry.
	// Return true if the entry is now locked, or false if the entry was already locked.
	Lock(logIndex raft.LogIndex, name string) bool

	// Unlock the given entry.
	// Return true if the entry is now unlocked, or false if the entry was already unlocked.
	Unlock(logIndex raft.LogIndex, name string) bool
}
