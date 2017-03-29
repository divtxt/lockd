package lockimpl

import (
	"log"

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
	var success bool
	if commitChan != nil {
		<-commitChan // FIXME: add timeout!
		success = true
	} else {
		success = false
	}
	log.Printf("lockd: Lock(\"%v\") -> %v", name, success)
	return &lockapi.LockResult{Success: success}, nil
}

func (s *LockingServerImpl) Unlock(ctx context.Context, in *lockapi.LockName) (*lockapi.LockResult, error) {
	name := in.Name
	commitChan := s.lockApi.Unlock(name)
	var success bool
	if commitChan != nil {
		<-commitChan // FIXME: add timeout!
		success = true
	} else {
		success = false
	}
	log.Printf("lockd: Unlock(\"%v\") -> %v", name, success)
	return &lockapi.LockResult{Success: success}, nil
}
