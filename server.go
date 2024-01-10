package gem

import (
	"errors"
	"sync"
)

func init() {
}

var (
	ErrSessionAlreadyExist = errors.New("session already exists")
)

type Server struct {
	id       int64
	sessions sync.Map
	handler  SessionHandler
	registry Registry
}

func NewServer(id int64) *Server {
	s := &Server{
		id:       id,
		sessions: sync.Map{},
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

// 加载session
// 如果会话已经存在，则直接返回
// 如果会话不存在，则创建，并调用SessionHandler.OnSessionLoad
// 如果开启注册表功能，则创建前要先到注册表查询状态，如果注册表中已经存在，则将旧值顶下线。
func (server *Server) LoadOrStoreSession(userID UserID) (Session, error) {
	if r := server.registry; r != nil {
		if err := r.AddUser(userID, &forceLogout{server}); err != nil {
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

// 处理客户端发过来的请求
func (server *Server) Request(userID UserID, req Request) {
	if server.handler == nil {
		return
	}
	if s, err := server.LoadOrStoreSession(userID); err != nil {
		return
	} else {
		server.handler.OnSessionRequest(s, req)
	}
}

// 注册表
func (server *Server) WithRegistry(registry Registry) {
	server.registry = registry
}
