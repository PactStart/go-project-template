package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/utils"
	"time"
)

const (
	TOKEN_PREFIX = "token:"
)

type TokenCache interface {
	AddTokenFlag(ctx context.Context, userID int64, platformID int, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID int64, platformID int) (map[string]int, error)
	SetTokenMapByUidPid(ctx context.Context, userID int64, platform int, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID int64, platform int, fields []string) error
}

type tokenCache struct {
	rdb        redis.UniversalClient
	expireTime time.Duration
}

func NewTokenCache(client redis.UniversalClient) TokenCache {
	return &tokenCache{rdb: client}
}

func (c *tokenCache) AddTokenFlag(ctx context.Context, userID int64, platformID int, token string, flag int) error {
	key := TOKEN_PREFIX + utils.Int64ToString(userID) + ":" + constant.PlatformIDToName(platformID)

	return errs.Wrap(c.rdb.HSet(ctx, key, token, flag).Err())
}

func (c *tokenCache) GetTokensWithoutError(ctx context.Context, userID int64, platformID int) (map[string]int, error) {
	key := TOKEN_PREFIX + utils.Int64ToString(userID) + ":" + constant.PlatformIDToName(platformID)
	m, err := c.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = utils.StringToInt(v)
	}

	return mm, nil
}

func (c *tokenCache) SetTokenMapByUidPid(ctx context.Context, userID int64, platform int, m map[string]int) error {
	key := TOKEN_PREFIX + utils.Int64ToString(userID) + ":" + constant.PlatformIDToName(platform)
	mm := make(map[string]interface{})
	for k, v := range m {
		mm[k] = v
	}

	return errs.Wrap(c.rdb.HSet(ctx, key, mm).Err())
}

func (c *tokenCache) DeleteTokenByUidPid(ctx context.Context, userID int64, platform int, fields []string) error {
	key := TOKEN_PREFIX + utils.Int64ToString(userID) + ":" + constant.PlatformIDToName(platform)

	return errs.Wrap(c.rdb.HDel(ctx, key, fields...).Err())
}
