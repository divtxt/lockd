package backend

import (
	"sync"

	"github.com/divtxt/lockd/util"
	"github.com/divtxt/raft"
)

// InMemoryBackend is an in-memory implementation of LockBackend.
type InMemoryBackend struct {
	mutex       sync.Mutex
	lastApplied raft.LogIndex
	locks       map[string]bool
}

func NewInMemoryBackend(initialLastApplied raft.LogIndex) *InMemoryBackend {
	return &InMemoryBackend{
		lastApplied: initialLastApplied,
		locks:       make(map[string]bool),
	}
}

// Get the current value of lastApplied.
func (imb *InMemoryBackend) GetLastApplied() raft.LogIndex {
	imb.mutex.Lock()
	defer imb.mutex.Unlock()

	return imb.lastApplied
}

// Check if the given entry is locked.
func (imb *InMemoryBackend) IsLocked(name string) bool {
	imb.mutex.Lock()
	defer imb.mutex.Unlock()

	if e := util.IsValidLockName(name); e != "" {
		panic(e)
	}
	_, ok := imb.locks[name]
	return ok
}

// Lock the given entry.
// Return true if the entry is now locked, or false if the entry was already locked.
func (imb *InMemoryBackend) Lock(logIndex raft.LogIndex, name string) bool {
	imb.mutex.Lock()
	defer imb.mutex.Unlock()

	if e := util.IsValidLockName(name); e != "" {
		panic(e)
	}
	// in-memory so we updating lastApplied is a trivial cost
	imb.lastApplied = logIndex
	// if already locked return false
	if _, ok := imb.locks[name]; ok {
		return false
	}
	// lock
	imb.locks[name] = true
	return true
}

// Unlock the given entry.
// Return true if the entry is now unlocked, or false if the entry was already unlocked.
func (imb *InMemoryBackend) Unlock(logIndex raft.LogIndex, name string) bool {
	imb.mutex.Lock()
	defer imb.mutex.Unlock()

	if e := util.IsValidLockName(name); e != "" {
		panic(e)
	}
	// in-memory so we updating lastApplied is a trivial cost
	imb.lastApplied = logIndex
	// if not locked return false
	if _, ok := imb.locks[name]; !ok {
		return false
	}
	// unlock
	delete(imb.locks, name)
	return true
}
