package components

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"orderin-server/pkg/common/application"
	"orderin-server/pkg/common/log"
)

const (
	CHANNEL_CONFIG_REFRESH = "config_refresh"
)

func Subscribe(rdb redis.UniversalClient, channelName string) {
	ctx := context.Background()
	// 订阅频道
	pubsub := rdb.Subscribe(ctx, channelName)
	channel := pubsub.Channel()
	go func() {
		log.ZInfo(ctx, "channel订阅成功", "channel", channelName)
		for msg := range channel {
			defer func() {
				if err := recover(); err != nil {
					log.ZError(ctx, "处理消息发生错误", errors.New("Server Error"), "err", err)
				}
			}()
			log.ZInfo(ctx, "channel收到一条消息", "channel", msg.Channel, "payload", msg.Payload)
			switch msg.Channel {
			case CHANNEL_CONFIG_REFRESH:
				application.AppContext.GetComponent(COMPONENT_MY_CONFIG).(*MyConfig).ReloadAll()
			}
		}
	}()
}

func Publish(rdb redis.UniversalClient, channelName string, message interface{}) error {
	ctx := context.Background()
	return rdb.Publish(ctx, channelName, message).Err()
}
