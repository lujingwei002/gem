package gem

type UserID interface {
	Value() any
}

type SessionId interface {
	Value() any
}

type Session interface {
	// 处理客户端的请求，并等待回复
	Request(req Request) Response
	UserID() UserID
	Server() *Server
}

// 会话事件回调
type SessionHandler interface {
	OnSessionCreated(s Session) (Session, error)
	OnSessionDestory(s Session)
	OnSessionForceLogout(s Session)
	OnSessionRequest(s Session, req Request)
}

type session struct {
	userID UserID
	server *Server
}

func (s *session) Request(req Request) Response {
	return nil
}

func (s *session) UserID() UserID {
	return s.userID
}

func (s *session) Server() *Server {
	return s.server
}

func newSession(server *Server, userID UserID) *session {
	s := &session{
		server: server,
		userID: userID,
	}
	return s
}
