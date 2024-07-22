package cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/errs"
	"time"
)

const (
	maxRetry = 10 // number of retries
)

// NewRedis Initialize redis connection.
func NewRedis() (redis.UniversalClient, error) {
	if len(config.Config.Redis.Address) == 0 {
		return nil, errors.New("redis address is empty")
	}
	errs.AddReplace(redis.Nil, errs.ErrRecordNotFound)
	var rdb redis.UniversalClient
	if len(config.Config.Redis.Address) > 1 || config.Config.Redis.ClusterMode {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:      config.Config.Redis.Address,
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password,
			PoolSize:   config.Config.Redis.PoolSize,
			MaxRetries: maxRetry,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:       config.Config.Redis.Address[0],
			Username:   config.Config.Redis.Username,
			Password:   config.Config.Redis.Password,
			DB:         config.Config.Redis.DB,
			PoolSize:   config.Config.Redis.PoolSize,
			MaxRetries: maxRetry,
		})
	}

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = rdb.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("redis ping %w", err)
	}

	return rdb, err
}
