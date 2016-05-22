package lockd

// Implementation of LockApi
func NewLockApiImpl(lsp LockStatePersistence) LockApi {
	return &LockApiImpl{lsp}
}

type LockApiImpl struct {
	lsp LockStatePersistence
}

func (l *LockApiImpl) Lock(name string) (bool, error) {
	return l.lsp.Lock(name)
}

func (l *LockApiImpl) Unlock(name string) (bool, error) {
	return l.lsp.Unlock(name)
}
