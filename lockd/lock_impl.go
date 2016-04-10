package lockd

// Implementation of LockApi
func NewLockApiImpl() LockApi {
	return &LockApiImpl{}
}

type LockApiImpl struct {
}

func (l *LockApiImpl) Lock(name string) (bool, error) {
	return false, nil
}

func (l *LockApiImpl) Unlock(name string) (bool, error) {
	return false, nil
}
