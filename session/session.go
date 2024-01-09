package session

import (
	"sync"

	"github.com/lujingwei002/gem"
)

func init() {
	server = New()
}

var (
	server *Server
)

type Id int64

func (id Id) Value() any {
	return id
}

type Server struct {
	sessions sync.Map
	handler  gem.SessionHandler
}

func New() *Server {
	s := &Server{
		sessions: sync.Map{},
	}
	return s
}

type Session struct {
}

func (s *Session) Request(req gem.Request) gem.Response {
	return nil
}

func (server *Server) RegisterGrpcResolver() error {
	return nil
}

// 加载session
// 如果会话已经存在，则直接返回
// 如果会话不存在，则创建，并调用SessionHandler.OnSessionLoad
// 如果开启注册表功能，则创建前要先到注册表查询状态，如果注册表中已经存在，则将旧值顶下线。
func (server *Server) Load(id gem.SessionId) (gem.Session, error) {
	if s, loaded := server.sessions.LoadOrStore(id.Value(), &Session{}); !loaded {
		if s, err := server.handler.OnSessionLoad(s.(gem.Session)); err != nil {
			return nil, err
		} else {
			return s, nil
		}
	} else {
		return s.(gem.Session), nil
	}
}

// 注册handler
func (server *Server) WithHandler(handler gem.SessionHandler) *Server {
	server.handler = handler
	return server
}

// 加载session
func Load(id gem.SessionId) (gem.Session, error) {
	return server.Load(id)
}

// 注册handler
func WithHandler(handler gem.SessionHandler) *Server {
	return server.WithHandler(handler)
}

// 注册成grpc的resolver
func RegisterGrpcResolver() error {
	return server.RegisterGrpcResolver()
}
