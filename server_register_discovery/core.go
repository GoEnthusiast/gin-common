package server_register_discovery

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoEnthusiast/gin-common/golangx/slicex"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"net"
	"time"
)

var ServerNameUndefined = errors.New("server name is undefined")
var RedisPoolUndefined = errors.New("redis pool is undefined")
var TimeoutError = errors.New("Expiration time is too short")

type (
	ServerRegisterDiscovery interface {
		Register(string) error
		Discover(string) (string, error)
		GetLocalAddr() (string, error)
	}
	serverRegisterDiscovery struct {
		redisPool        *redis.Client
		serverName       string
		serverAddTimeout int
		//db               int
	}
	Option func(encry *serverRegisterDiscovery)
)

func (s serverRegisterDiscovery) GetLocalAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	var localIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				break
			}
		}
	}
	return localIP, nil
}

func (s serverRegisterDiscovery) Register(addr string) error {
	if s.serverName == ":" {
		return ServerNameUndefined
	}
	if s.serverAddTimeout < 1 {
		return TimeoutError
	}
	ctx := context.Background()
	s.redisPool.Set(ctx, s.serverName, addr, time.Duration(s.serverAddTimeout)*time.Second)
	return nil
}

func (s serverRegisterDiscovery) Discover(serverName string) (string, error) {
	ctx := context.Background()
	var cursor uint64
	var n int
	match := fmt.Sprintf("%s:*", serverName)
	var result []string
	for {
		var keys []string
		var err error
		keys, cursor, err = s.redisPool.Scan(ctx, cursor, match, 10000).Result()
		if err != nil {
			return "", err
		}
		n += len(keys)
		if len(keys) > 0 {
			result = append(result, keys...)
		}
		if cursor == 0 {
			break
		}
	}

	serverKey, err := slicex.GetSliceValueWithRand(result)
	if err != nil {
		return "", err
	}
	serverKeyStr := serverKey.(string)
	addr := s.redisPool.Get(ctx, serverKeyStr).Val()
	if addr == "" {
		return "", errors.New("not found addr")
	}
	return addr, nil
}

func NewServerRegisterDescovery(opts ...Option) (ServerRegisterDiscovery, error) {
	var s serverRegisterDiscovery
	for _, opt := range opts {
		opt(&s)
	}
	if s.serverAddTimeout == 0 {
		s.serverAddTimeout = 10
	}
	if s.serverName == "" {
		s.serverName = ":"
	} else {
		s.serverName = fmt.Sprintf("%s:%s%v", s.serverName, uuid.New().String(), time.Now().UnixNano())
	}
	return &s, nil
}

func WithServerName(arg string) Option {
	return func(encry *serverRegisterDiscovery) {
		encry.serverName = arg
	}
}

func WithRedisPool(arg *redis.Client) Option {
	return func(encry *serverRegisterDiscovery) {
		encry.redisPool = arg
	}
}

func WithServerAddTimeout(arg int) Option {
	return func(encry *serverRegisterDiscovery) {
		encry.serverAddTimeout = arg
	}
}

//func WithDB(arg int) Option {
//	return func(encry *serverRegisterDiscovery) {
//		encry.db = arg
//	}
//}
