package gem

type SessionId interface {
	Value() any
}

type Request interface {
}

type Response interface {
}

type Session interface {
	// 处理客户端的请求，并等待回复
	Request(req Request) Response
}

// 会话事件回调
type SessionHandler interface {
	OnSessionLoad(s Session) (Session, error)
	OnSessionDestory(s Session)
}
