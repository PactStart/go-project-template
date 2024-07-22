package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"time"
)

type VerifyCodeCache interface {
	StoreCode(id string, scene string, code string, expireTime time.Duration) error
	CheckCode(context context.Context, id string, scene string, code string) (bool, error)
}

type verifyCodeCache struct {
	rdb redis.UniversalClient
}

func (s verifyCodeCache) StoreCode(id string, scene string, code string, expireTime time.Duration) error {
	key := "code:" + scene + ":" + id
	return s.rdb.Set(context.Background(), key, code, expireTime).Err()
}

func (s verifyCodeCache) CheckCode(context context.Context, id string, scene string, code string) (bool, error) {
	key := "code:" + scene + ":" + id
	exists, err := s.rdb.Exists(context, key).Result()
	if err != nil {
		return false, errs.NewCodeError(errs.SmsCodeError, "验证码已过期").WithDetail(err.Error())
	}
	if exists == 1 {
		result, err := s.rdb.Get(context, key).Result()
		if err != nil {
			return false, errs.NewCodeError(errs.SmsCodeError, "验证码错误").WithDetail(err.Error())
		}
		if result != code {
			log.ZError(context, "验证码错误", err, "expect", result, "actual", code)
			return false, errs.NewCodeError(errs.SmsCodeError, "验证码错误")
		}
		return true, nil
	} else {
		return false, errs.NewCodeError(errs.SmsCodeError, "验证码已过期")
	}
}

func NewVerifyCodeCache(client redis.UniversalClient) VerifyCodeCache {
	return &verifyCodeCache{rdb: client}
}
