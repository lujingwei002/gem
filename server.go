package gem

import (
	"context"
	"sync"

	"github.com/lujingwei002/gem/services/registrypb"
	"google.golang.org/grpc"
)

func init() {
}

var ()

type Server struct {
	id             int64
	sessions       sync.Map
	handler        SessionHandler
	registry       Registry
	address        string // 内部通信地址
	gRPC           *grpc.Server
	registryServer *registry_server
}

func NewServer(id int64) *Server {
	s := &Server{
		id:             id,
		sessions:       sync.Map{},
		registryServer: &registry_server{},
	}
	return s
}

type forceLogout struct {
	*Server
}

func (server *forceLogout) ForceLogout(userID UserID) error {
	if s, loaded := server.sessions.LoadAndDelete(userID.Value()); loaded {
		if server.handler != nil {
			server.handler.OnSessionForceLogout(s.(Session))
		}
	}
	return nil
}

func (server *Server) RegisterGrpcResolver() error {
	return nil
}

func (server *Server) ServerID() int64 {
	return server.id
}

func (server *Server) ServerAddress() string {
	return server.address
}

// 加载session
// 如果会话已经存在，则直接返回
// 如果会话不存在，则创建，并调用SessionHandler.OnSessionLoad
// 如果开启注册表功能，则创建前要先到注册表查询状态，如果注册表中已经存在，则将旧值顶下线。
func (server *Server) LoadOrStoreSession(ctx context.Context, userID UserID) (Session, error) {
	if r := server.registry; r != nil {
		if err := r.AddUser(ctx, userID, &forceLogout{server}); err != nil {
			return nil, err
		}
		if s, loaded := server.sessions.LoadOrStore(userID.Value(), newSession(server, userID)); !loaded {
			if s, err := server.handler.OnSessionCreated(s.(Session)); err != nil {
				return nil, err
			} else {
				return s, nil
			}
		} else {
			return s.(Session), nil
		}
	} else {
		if s, loaded := server.sessions.LoadOrStore(userID.Value(), newSession(server, userID)); !loaded {
			if s, err := server.handler.OnSessionCreated(s.(Session)); err != nil {
				return nil, err
			} else {
				return s, nil
			}
		} else {
			return s.(Session), nil

		}
	}
}

// 注册handler
func (server *Server) WithHandler(handler SessionHandler) *Server {
	server.handler = handler
	return server
}

// 内部通信地址
func (server *Server) WithAddress(address string) *Server {
	server.address = address
	return server
}

// 处理客户端发过来的请求
// 此操作会阻塞，直接得到一个response或者出错
func (server *Server) Request(ctx context.Context, userID UserID, req Request) {
	if server.handler == nil {
		return
	}
	if s, err := server.LoadOrStoreSession(ctx, userID); err != nil {
		return
	} else {
		server.handler.OnSessionRequest(s, req)
	}
}

// 注册表
func (server *Server) WithRegistry(registry Registry) {
	server.registry = registry
}

// 集成grpc
func (server *Server) WithGrpcServer(s *grpc.Server) {
	server.gRPC = s
	registrypb.RegisterRegistryServer(s, server.registryServer)
}
