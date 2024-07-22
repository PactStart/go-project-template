package components

import (
	"context"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"orderin-server/pkg/common/utils"
	"time"
)

type MyCache struct {
	rdb redis.UniversalClient
}

func (e *MyCache) SetSkipCheckIdentity(id int64, operation string) error {
	key := "skip_check_identity:" + utils.Int64ToString(id) + ":" + operation
	return e.rdb.Set(context.Background(), key, "", 30*time.Minute).Err()
}

func (e *MyCache) IsCheckedIdentity(id int64, operation string) (bool, error) {
	key := "skip_check_identity:" + utils.Int64ToString(id) + ":" + operation
	count, err := e.rdb.Exists(context.Background(), key).Result()
	return count > 0, err
}

func (e *MyCache) SetMyPermissions(id int64, permissons []string) error {
	key := "acl:" + utils.Int64ToString(id)
	e.rdb.Del(context.Background(), key)
	err := e.rdb.SAdd(context.Background(), key, permissons).Err()
	return err
}

func (e *MyCache) HasPerm(id int64, perm string) bool {
	key := "acl:" + utils.Int64ToString(id)
	exists, err := e.rdb.SIsMember(context.Background(), key, perm).Result()
	if err != nil {
		return false
	}
	return exists
}

func (e *MyCache) SetWechatLoginResult(state string, id int64) {
	key := "wechat_login:" + state
	e.rdb.Set(context.Background(), key, utils.Int64ToString(id), 10*time.Minute)
}

func (e *MyCache) GetWechatLoginResult(state string) (string, error) {
	key := "wechat_login:" + state
	result, err := e.rdb.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	return result, err

}
