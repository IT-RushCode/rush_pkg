package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/IT-RushCode/rush_pkg/config"
)

func REDIS_CONNECT(ctx context.Context, cfg *config.RedisConfig, modifiers ...func(*redis.Options)) *redis.Client {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.HOST, cfg.PORT),
		Password: cfg.PASS,
		DB:       cfg.DB,
		Protocol: 3,
	}

	if cfg.MaxRetries > 0 {
		opts.MaxRetries = cfg.MaxRetries
	}
	if cfg.MinRetryBackoff > 0 {
		opts.MinRetryBackoff = cfg.MinRetryBackoff
	}
	if cfg.MaxRetryBackoff > 0 {
		opts.MaxRetryBackoff = cfg.MaxRetryBackoff
	}
	if cfg.DialTimeout > 0 {
		opts.DialTimeout = cfg.DialTimeout
	}
	if cfg.DialerRetries > 0 {
		opts.DialerRetries = cfg.DialerRetries
	}
	if cfg.DialerRetryTime > 0 {
		opts.DialerRetryTimeout = cfg.DialerRetryTime
	}
	if cfg.ReadTimeout > 0 {
		opts.ReadTimeout = cfg.ReadTimeout
	}
	if cfg.WriteTimeout > 0 {
		opts.WriteTimeout = cfg.WriteTimeout
	}
	if cfg.ContextTimeout {
		opts.ContextTimeoutEnabled = true
	}
	if cfg.PoolSize > 0 {
		opts.PoolSize = cfg.PoolSize
	}
	if cfg.MaxActiveConns > 0 {
		opts.MaxActiveConns = cfg.MaxActiveConns
	}
	if cfg.MaxConcurrentDials > 0 {
		opts.MaxConcurrentDials = cfg.MaxConcurrentDials
	}
	if cfg.PoolTimeout > 0 {
		opts.PoolTimeout = cfg.PoolTimeout
	}
	if cfg.MinIdleConns > 0 {
		opts.MinIdleConns = cfg.MinIdleConns
	}
	if cfg.MaxIdleConns > 0 {
		opts.MaxIdleConns = cfg.MaxIdleConns
	}
	if cfg.ConnMaxIdleTime > 0 {
		opts.ConnMaxIdleTime = cfg.ConnMaxIdleTime
	}
	if cfg.ConnMaxLifetime > 0 {
		opts.ConnMaxLifetime = cfg.ConnMaxLifetime
	}

	for _, mod := range modifiers {
		if mod != nil {
			mod(opts)
		}
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		panic("redis is not connected: " + err.Error())
	}

	return client
}

func REDIS_CLOSE(client *redis.Client) {
	_ = client.Close()
}
