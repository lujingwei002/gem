package session

import (
	"fmt"
	"testing"

	"github.com/lujingwei002/gem"
	"github.com/lujingwei002/gem/request"
)

type TestSessionHandler struct {
}

func (TestSessionHandler) OnSessionLoad(s gem.Session) (gem.Session, error) {
	fmt.Println("Load")
	return s, nil
}

func (TestSessionHandler) OnSessionDestory(s gem.Session) {
	fmt.Println("Destory")
}

func TestSessionLoad(t *testing.T) {
	server := New()
	server.WithHandler(&TestSessionHandler{})
	// 注册grpc resolver
	server.RegisterGrpcResolver()
	// 会话注册到注册表

	id := Id(1)
	s, _ := server.Load(id)
	req := request.Request{}
	s.Request(req)
}
