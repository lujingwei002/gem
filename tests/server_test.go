package tests

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/lujingwei002/gem"
	"github.com/lujingwei002/gem/proto/dialog"
	"github.com/lujingwei002/gem/registry/local_registry"
	"github.com/lujingwei002/gem/registry/redis_registry"
	"github.com/lujingwei002/gem/userid"
	"google.golang.org/grpc"
)

type TestSessionHandler struct {
	onSessionCreated     func(s gem.Session)
	onSessionRequest     func(s gem.Session)
	onSessionForceLogout func(s gem.Session)
}

func (h TestSessionHandler) OnSessionCreated(s gem.Session) (gem.Session, error) {
	// slog.Info("OnSessionCreated", "server_id", s.Server().ServerID(), "user_id", s.UserID())
	if h.onSessionCreated != nil {
		h.onSessionCreated(s)
	}
	return s, nil
}

func (h TestSessionHandler) OnSessionDestory(s gem.Session) {
	// slog.Info("OnSessionDestory", "server_id", s.Server().ServerID(), "user_id", s.UserID())
}

func (h TestSessionHandler) OnSessionRequest(s gem.Session, req gem.Request) {
	// slog.Info("OnSessionRequest", "server_id", s.Server().ServerID(), "user_id", s.UserID())
	if h.onSessionRequest != nil {
		h.onSessionRequest(s)
	}
}

func (h TestSessionHandler) OnSessionForceLogout(s gem.Session) {
	// slog.Info("OnSessionForceLogout", "server_id", s.Server().ServerID(), "user_id", s.UserID())
	if h.onSessionForceLogout != nil {
		h.onSessionForceLogout(s)
	}
}

func TestSessionLoad(t *testing.T) {

	server1 := gem.NewServer(1)
	createdTimes := 0
	requestTimes := 0
	server1.WithHandler(&TestSessionHandler{
		onSessionCreated: func(s gem.Session) {
			createdTimes++
		},
		onSessionRequest: func(s gem.Session) {
			requestTimes++
		},
	})

	userId := userid.Int64(1)
	req := dialog.New("ljw", "hello")

	var times int = 4
	for i := 1; i <= times; i++ {
		server1.Request(context.Background(), userId, req)
	}
	if createdTimes != 1 {
		t.Fatalf("session created times failed, expected=%d got=%d", 1, createdTimes)
	}
	if requestTimes != times {
		t.Fatalf("session request times failed, expected=%d got=%d", times, requestTimes)
	}
}

func TestSessionLocalForceLogout(t *testing.T) {
	server1 := gem.NewServer(1)
	server1.WithHandler(&TestSessionHandler{})

	server2 := gem.NewServer(2)
	server2.WithHandler(&TestSessionHandler{})
	// 注册grpc resolver
	// server.RegisterGrpcResolver()
	// 会话注册到注册表
	reg := local_registry.New()
	server1.WithRegistry(reg)
	server2.WithRegistry(reg)

	userId := userid.Int64(1)
	req := dialog.New("ljw", "hello")

	server1.Request(context.Background(), userId, req)
	server2.Request(context.Background(), userId, req)
}

func TestSessionRedisRegistry(t *testing.T) {
	grpcServer1 := grpc.NewServer()
	server1 := gem.NewServer(1)
	server1.WithHandler(&TestSessionHandler{}).
		WithAddress("127.0.0.1:4441").
		WithGrpcServer(grpcServer1)

	grpcServer2 := grpc.NewServer()
	server2 := gem.NewServer(2)
	server2.WithHandler(&TestSessionHandler{}).
		WithAddress("127.0.0.1:4442").
		WithGrpcServer(grpcServer2)

	// 注册grpc resolver
	// server.RegisterGrpcResolver()
	// 会话注册到注册表
	reg1, _ := redis_registry.Connect(context.Background(), "127.0.0.1:6379", "123456", 0)
	server1.WithRegistry(reg1)

	reg2, _ := redis_registry.Connect(context.Background(), "127.0.0.1:6379", "123456", 0)
	server2.WithRegistry(reg2)

	userId := userid.Int64(1)
	req := dialog.New("ljw", "hello")

	server1.Request(context.Background(), userId, req)
	server2.Request(context.Background(), userId, req)
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	var urls = []string{
		"http://www.baidu.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			http.Get(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
	fmt.Println("eeeeeeeeeee1")

	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			http.Get(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
	fmt.Println("eeeeeeeeeee2")
}
