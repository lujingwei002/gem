package local_registry

import (
	"context"
	"errors"
	"sync"

	"github.com/lujingwei002/gem"
)

var (
	ErrForceLogoutUndefined = errors.New("ForceLogout undefined")
	ErrForceLogout          = errors.New("ForceLogout failed")
)

type LocalRegistry struct {
	users sync.Map
}

func New() *LocalRegistry {
	self := &LocalRegistry{
		users: sync.Map{},
	}
	return self
}

// 添加用户到注册表，如果用户之前已经存在，则将顶号顶下线
func (r *LocalRegistry) AddUser(ctx context.Context, userID gem.UserID, new gem.RegistryServer) error {
	if new == nil {
		return ErrForceLogoutUndefined
	}
	if old, loaded := r.users.LoadOrStore(userID.Value(), new); loaded {
		if err := old.(gem.RegistryServer).ForceLogout(userID); err != nil {
			return err
		}
		if ok := r.users.CompareAndSwap(userID.Value(), old, new); !ok {
			return ErrForceLogout
		}
		return nil
	} else {
		return nil
	}
}
