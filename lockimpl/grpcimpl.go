package lockimpl

import (
	"golang.org/x/net/context"

	"github.com/divtxt/lockd/lockapi"
)

func NewLockingServerImpl(lockApi InternalLockApi) *LockingServerImpl {
	return &LockingServerImpl{lockApi}
}

type LockingServerImpl struct {
	lockApi InternalLockApi
}

func (s *LockingServerImpl) IsLocked(ctx context.Context, in *lockapi.LockName) (*lockapi.LockState, error) {
	name := in.Name
	locked, _ := s.lockApi.IsLocked(name)
	return &lockapi.LockState{Name: name, IsLocked: locked}, nil
}

func (s *LockingServerImpl) Lock(ctx context.Context, in *lockapi.LockName) (*lockapi.LockResult, error) {
	name := in.Name
	commitChan := s.lockApi.Lock(name)
	if commitChan != nil {
		<-commitChan // FIXME: add timeout!
		return &lockapi.LockResult{Success: true}, nil
	} else {
		return &lockapi.LockResult{Success: false}, nil
	}
}

func (s *LockingServerImpl) Unlock(ctx context.Context, in *lockapi.LockName) (*lockapi.LockResult, error) {
	name := in.Name
	commitChan := s.lockApi.Unlock(name)
	if commitChan != nil {
		<-commitChan // FIXME: add timeout!
		return &lockapi.LockResult{Success: true}, nil
	} else {
		return &lockapi.LockResult{Success: false}, nil
	}
}
