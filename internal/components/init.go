package components

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/application"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/weixin"
)

const (
	COMPONENT_ORM       = "orm"
	COMPONENT_REDIS     = "redis"
	COMPONENT_MY_CACHE  = "my_cache"
	COMPONENT_MY_CONFIG = "my_config"
)

func RegisterComponents(db *gorm.DB, rdb redis.UniversalClient) {
	application.NewApplicationContext()

	application.AppContext.RegisterComponent(COMPONENT_ORM, db)
	log.ZInfo(context.Background(), "orm注册成功")

	application.AppContext.RegisterComponent(COMPONENT_REDIS, rdb)
	log.ZInfo(context.Background(), "redis注册成功")

	myCache := MyCache{rdb: rdb}
	application.AppContext.RegisterComponent(COMPONENT_MY_CACHE, &myCache)
	log.ZInfo(context.Background(), "my_cache注册成功")

	configService := services.SysConfig{}
	configService.Orm = db
	configService.Context = context.Background()
	myConfig := MyConfig{ConfigService: configService, ConfigMap: make(map[string]string)}
	myConfig.ReloadAll()
	application.AppContext.RegisterComponent(COMPONENT_MY_CONFIG, &myConfig)
	log.ZInfo(context.Background(), "my_config注册成功")

}

func RegisterAuthorizers(db *gorm.DB) {
	if config.Config.Env.Profiles == "dev" {
		return
	}

	context := context.Background()
	authorizerService := services.WxAuthorizer{}
	authorizerService.Orm = db
	authorizerService.Context = context

	authorizers, err := authorizerService.GetAll()
	if err != nil {
		log.ZError(authorizerService.Context, "query authorizers fail", err)
		return
	}
	for _, authorizer := range *authorizers {
		RegisterAuthorizer(context, authorizer)
	}
}

func RegisterAuthorizer(context context.Context, authorizer models.WxAuthorizer) {
	officialAccount, err := weixin.OpenPlatformApp.OfficialAccount(authorizer.AuthorizerAppid, authorizer.AuthorizerRefreshToken, nil)
	if err == nil {
		application.AppContext.RegisterComponent(authorizer.AuthorizerAppid, officialAccount)
		log.ZInfo(context, "微信公众号注册成功", "appId", authorizer.AuthorizerAppid)
	} else {
		log.ZError(context, "build authorizer officialAccount fail", err, "appId", authorizer.AuthorizerAppid, "refreshToken", authorizer.AuthorizerRefreshToken)
	}
}
