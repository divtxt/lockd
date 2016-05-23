package lockd

// In-memory implementation of LockStatePersistence
func NewInMemoryLSP() *InMemoryLSP {
	return &InMemoryLSP{
		make(map[string]bool),
	}
}

type InMemoryLSP struct {
	locks map[string]bool
}

func (l *InMemoryLSP) Lock(name string) (bool, error) {
	if _, ok := l.locks[name]; ok {
		return false, nil
	}
	l.locks[name] = true
	return true, nil
}

func (l *InMemoryLSP) Unlock(name string) (bool, error) {
	if _, ok := l.locks[name]; ok {
		delete(l.locks, name)
		return true, nil
	}
	return false, nil
}
