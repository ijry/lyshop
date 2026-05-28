package cache

import (
	"context"
	"fmt"
	"strings"

	"github.com/alicebob/miniredis/v2"
	"github.com/ijry/lyshop/config"
	"github.com/redis/go-redis/v9"
)

// Client is the global Redis client, available after Init().
var Client *redis.Client
var embedded *miniredis.Miniredis

func Init() error {
	cfg := config.Global.Redis
	addr := strings.TrimSpace(cfg.Addr)
	password := cfg.Password
	if addr == "" {
		srv, err := miniredis.Run()
		if err != nil {
			return fmt.Errorf("start embedded redis: %w", err)
		}
		embedded = srv
		addr = srv.Addr()
		password = ""
	}

	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       cfg.DB,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		if embedded != nil {
			embedded.Close()
			embedded = nil
		}
		return fmt.Errorf("redis connect: %w", err)
	}
	Client = c
	return nil
}
