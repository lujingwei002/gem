package redis_registry

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lujingwei002/gem"
)

const (
	PREFIX = "user_location:"
)

var (
	ErrForceLogout = errors.New("ForceLogout failed")
	ErrDeny        = errors.New("deny, it's not your user")
)

// 脚本
const (
	// KEYS[1]=key  ARGS[1]=旧值  返回 {删除的key数量,旧值}
	script_compare_and_delete = `
local v = redis.call("GET", KEYS[1])
if v == false then
	return {0, v}
elseif v ~= ARGV[1] then
	return {0, v}
else
	return {redis.call("DEL", KEYS[1]), v}
end
`
	// KEYS[1]=key ARGS[1]=值 ARGS[3]=expiration 返回
	script_setnx_and_get = `
local v = redis.call("GET", KEYS[1])
if redis.call("SETNX", KEYS[1], ARGV[1]) == 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[2])
	return {1, v}
else
	return {0, v}
end
`
)

type RedisRegistry struct {
	rdb        *redis.Client
	expiration time.Duration
	prefix     string
	address    string
	users      map[gem.UserID]*user
	mu         sync.Mutex

	compareAndDeleteScript string
	setNXAndGetScript      string
}

type user struct {
	status int
	server gem.RegistryServer
	wg     sync.WaitGroup
}

const (
	user_added    = 1
	user_adding   = 2
	user_removing = 3
)

// 将script字符串转换为sha
func loadScriptSha(ctx context.Context, rdb *redis.Client, script string) (string, error) {
	s := redis.NewScript(script)
	if result := s.Load(ctx, rdb); result.Err() != nil {
		return "", result.Err()
	} else {
		return result.Val(), nil
	}
}

func Connect(ctx context.Context, address string, password string, db int) (*RedisRegistry, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	compareAndDeleteScript, err := loadScriptSha(ctx, rdb, script_compare_and_delete)
	if err != nil {
		rdb.Close()
		return nil, err
	}

	setNXAndGetScript, err := loadScriptSha(ctx, rdb, script_setnx_and_get)
	if err != nil {
		rdb.Close()
		return nil, err
	}
	r := &RedisRegistry{
		users:                  make(map[gem.UserID]*user),
		rdb:                    rdb,
		address:                address,
		prefix:                 PREFIX,
		expiration:             1 * time.Second,
		mu:                     sync.Mutex{},
		compareAndDeleteScript: compareAndDeleteScript,
		setNXAndGetScript:      setNXAndGetScript,
	}
	return r, nil
}

func (r *RedisRegistry) fmtUserKey(userID gem.UserID) string {
	return fmt.Sprintf("%s%s", r.prefix, userID)
}

func (r *RedisRegistry) fmtServerValue(server gem.RegistryServer) string {
	return server.ServerAddress()
}

func (r *RedisRegistry) RemoveUser(ctx context.Context, userID gem.UserID, server gem.RegistryServer) error {
	r.mu.Lock()
	if user, ok := r.users[userID]; !ok {
		return nil
	} else if user.status == user_added && user.server == server {
		user.status = user_removing
		r.mu.Unlock()
		rkey := r.fmtUserKey(userID)
		if result := r.rdb.EvalSha(ctx, r.compareAndDeleteScript, []string{rkey}, r.fmtServerValue(server)); result.Err() != nil {
			user.status = user_added
			return result.Err()
		}
		r.mu.Lock()
		delete(r.users, userID)
		return nil
	} else if user.status == user_added && user.server != server {
		return ErrDeny
	} else if user.status == user_adding {
		r.mu.Unlock()
		user.wg.Wait()
	}
	return nil
}

func (r *RedisRegistry) AddUser(ctx context.Context, userID gem.UserID, server gem.RegistryServer) error {
	rkey := r.fmtUserKey(userID)
	//slog.Info("add user", "server_id", new.ServerID(), "server_address", new.ServerAddress(), "user_id", userID, "key", rkey)
	r.mu.Lock()
	if old, ok := r.users[userID]; !ok {
		new := &user{
			status: user_adding,
			server: server,
			wg:     sync.WaitGroup{},
		}
		new.wg.Add(1)
		r.users[userID] = new
		r.mu.Unlock()
		if result := r.rdb.EvalSha(ctx, r.setNXAndGetScript, []string{rkey}, r.fmtServerValue(server), r.expiration); result.Err() != nil {
			return result.Err()
		} else if _, err := result.Slice(); err != nil {
			new.wg.Done()
			return err
		}
		return nil
	} else if old.server == server {
		return nil
	} else if old.server != server {
		old.status = user_adding
		if result := r.rdb.EvalSha(ctx, r.setNXAndGetScript, []string{rkey}, r.fmtServerValue(server), r.expiration); result.Err() != nil {
			return result.Err()
		} else {
			result.Slice()
		}
	}
	return nil
}
