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

func (iml *InMemoryLSM) IsLocked(name string) (bool, error) {
	_, hasKey := iml.locks[name]
	return hasKey, nil
}

func (iml *InMemoryLSM) Lock(name string) (bool, error) {
	// if already locked return false
	if _, hasKey := iml.locks[name]; hasKey {
		return false, nil
	}
	// lock
	iml.locks[name] = true
	return true, nil
}

func (iml *InMemoryLSM) Unlock(name string) (bool, error) {
	// if not locked return false
	if _, hasKey := iml.locks[name]; !hasKey {
		return false, nil
	}
	// unlock
	delete(iml.locks, name)
	return true, nil
}
