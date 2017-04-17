package statemachine

// An in-memory implementation of LockStateMachine.
//
// This implementation is NOT concurrent safe.
//
type InMemoryLSM struct {
	locks map[string]bool
}

// Construct a new InMemoryLSM.
func NewInMemoryLSM() *InMemoryLSM {
	return &InMemoryLSM{make(map[string]bool)}
}

// --- Implement LockStateMachine

func (iml *InMemoryLSM) IsLocked(name string) bool {
	if e := IsValidLockName(name); e != "" {
		panic(e)
	}
	_, hasKey := iml.locks[name]
	return hasKey
}

func (iml *InMemoryLSM) Lock(name string) bool {
	if e := IsValidLockName(name); e != "" {
		panic(e)
	}
	// if already locked return false
	if _, hasKey := iml.locks[name]; hasKey {
		return false
	}
	// lock
	iml.locks[name] = true
	return true
}

func (iml *InMemoryLSM) Unlock(name string) bool {
	if e := IsValidLockName(name); e != "" {
		panic(e)
	}
	// if not locked return false
	if _, hasKey := iml.locks[name]; !hasKey {
		return false
	}
	// unlock
	delete(iml.locks, name)
	return true
}
