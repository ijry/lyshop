package cache

import (
	"context"
	"fmt"

	"github.com/ijry/lyshop/config"
	"github.com/redis/go-redis/v9"
)

// Client is the global Redis client, available after Init().
var Client *redis.Client

func Init() error {
	cfg := config.Global.Redis
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := c.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("redis connect: %w", err)
	}
	Client = c
	return nil
}
