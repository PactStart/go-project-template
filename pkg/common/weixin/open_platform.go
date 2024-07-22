package weixin

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/openPlatform"
	"orderin-server/pkg/common/config"
)

var OpenPlatformApp *openPlatform.OpenPlatform

func NewOpenPlatformAppService() (*openPlatform.OpenPlatform, error) {

	var cache kernel.CacheInterface

	if len(config.Config.Redis.Address) > 0 {
		cache = kernel.NewRedisClient(&kernel.UniversalOptions{
			Addrs:    config.Config.Redis.Address,
			Username: config.Config.Redis.Username,
			Password: config.Config.Redis.Password,
			DB:       config.Config.Redis.DB,
			PoolSize: config.Config.Redis.PoolSize,
		})
	}

	app, err := openPlatform.NewOpenPlatform(&openPlatform.UserConfig{

		AppID:  config.Config.WxOpenPlatform.AppID,
		Secret: config.Config.WxOpenPlatform.AppSecret,

		Token:  config.Config.WxOpenPlatform.MessageToken,
		AESKey: config.Config.WxOpenPlatform.MessageAesKey,

		Log: openPlatform.Log{
			Level: "debug",
			File:  "./logs/weixin.log",
		},
		Cache:     cache,
		HttpDebug: true,
		Debug:     true,
		//"sandbox": true,
	})

	return app, err
}
