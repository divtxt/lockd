package lockd

// In-memory implementation of LockStatePersistence
func NewInMemoryLSP() LockStatePersistence {
	return &InMemoryLSP{}
}

type InMemoryLSP struct {
}

func (l *InMemoryLSP) Lock(name string) (bool, error) {
	// FIXME: actually track lock state
	return false, nil
}

func (l *InMemoryLSP) Unlock(name string) (bool, error) {
	// FIXME: actually track lock state
	return false, nil
}
