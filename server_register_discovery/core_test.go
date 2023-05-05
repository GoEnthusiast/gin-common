package server_register_discovery

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

const serverName = "test-server"

func TestServerRegisterDiscovery_Register(t *testing.T) {
	// redis
	redisPool := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password: "",
		DB:       2,
		PoolSize: 100,
	})

	serverTimeout := 100
	for i := 1; i <= 5; i++ {
		s, err := NewServerRegisterDescovery(
			WithServerName(serverName),
			WithServerAddTimeout(serverTimeout),
			WithRedisPool(redisPool),
		)
		if err != nil {
			panic(err)
		}

		err = s.Register(fmt.Sprintf("127.0.0.%d", i))
		if err != nil {
			panic(err)
		}
	}
}

func TestServerRegisterDiscovery_Discover(t *testing.T) {
	// redis
	redisPool := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password: "",
		DB:       2,
		PoolSize: 100,
	})

	s, err := NewServerRegisterDescovery(
		WithRedisPool(redisPool),
	)
	if err != nil {
		panic(err)
	}

	redisKey, err := s.Discover(serverName)
	fmt.Println(redisKey)
}

func TestServerRegisterDiscovery_GetLocalAddr(t *testing.T) {
	s, err := NewServerRegisterDescovery()
	if err != nil {
		panic(err)
	}
	addr, err := s.GetLocalAddr()
	if err != nil {
		panic(err)
	}
	fmt.Println(addr)
}
