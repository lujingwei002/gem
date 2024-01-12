package gem

import (
	"context"

	"github.com/lujingwei002/gem/services/registrypb"
)

type Registry interface {
	// 注册账号
	// 如果账号已经存在，可能会发生顶号处理，由具体实现决定是否要顶号
	// 一旦此方法返回成功，如果需要顶号的话，则保证对方已经下线，且数据已经保存。
	AddUser(ctx context.Context, userID UserID, server RegistryServer) error
}

// 顶号下线
type RegistryServer interface {
	ForceLogout(userID UserID) error
	// 服务地址，用来设到注册表中，目前只能是grpc地址
	ServerAddress() string
	ServerID() int64
}

func (r *registry_server) ServerID() int64 {
	return r.server.id
}

func (r *registry_server) ServerAddress() string {
	return r.server.address
}

type registry_server struct {
	registrypb.UnimplementedRegistryServer
	server *Server
}

func (r *registry_server) ForceLogout(context.Context, *registrypb.ForceLogoutRequest) (*registrypb.ForceLogoutResponse, error) {
	resp := &registrypb.ForceLogoutResponse{}
	return resp, nil
}
