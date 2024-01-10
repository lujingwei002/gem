# gem
游戏服务端开发

## 插件列表
gem-resource


## Examples

启动服务,接收请求 

```go

server1 := gem.NewServer(1)

userId := userid.Int64(1)
req := dialog.New("ljw", "hello")

server1.Request(userId, req)
```

顶号下线，使用local_registry 

server1和server2需要可以运行在同一个的进程中 

```go
server1 := gem.NewServer(1)
server2 := gem.NewServer(2)

reg := local_registry.New()
server1.WithRegistry(reg)
server2.WithRegistry(reg)

userId := userid.Int64(1)
req := dialog.New("ljw", "hello")

server1.Request(userId, req)
// server1上的用户1会被顶下线，然后server2处理请求
server2.Request(userId, req)
```

顶号下线，使用redis_registry 

server1和server2可以运行在不同的进程中 

```go
server1 := gem.NewServer(1)
server2 := gem.NewServer(2)

reg := redis_registry.New()
server1.WithRegistry(reg)
server2.WithRegistry(reg)

userId := userid.Int64(1)
req := dialog.New("ljw", "hello")

server1.Request(userId, req)
// server1上的用户1会被顶下线，然后server2处理请求
server2.Request(userId, req)
```

从管道拉取请求 

```go
server1 := gem.NewServer(1)

userId := userid.Int64(1)
req := dialog.New("ljw", "hello")

p, c := channel_source.mpsc()

server1.PullFrom(c)

p <= req
p <= req

close(p)
```

从redis拉取请求 

```go
server1 := gem.NewServer(1)

userId := userid.Int64(1)
req := dialog.New("ljw", "hello")

source := redis_source.New()

server1.PullFrom(source)

source.Request(userId, req)
source.Request(userId, req)
```