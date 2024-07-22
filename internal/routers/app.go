package routers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
	docs "orderin-server/docs/app-api"
	"orderin-server/internal/api/app"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/ginmiddleware"
	"orderin-server/pkg/common/log"
)

// @title APP API
// @version 1.0.0
// @description xxx APP API文档

// @contact.name xxx
// @contact.url http://xxx.com
// @contact.email xxx@qq.com

// @host 127.0.0.1:10002
// @BasePath /api/v1
func NewH5GinRouter(db *gorm.DB, rdb redis.UniversalClient) *gin.Engine {
	if config.Config.Env.Profiles == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	log.ZInfo(context.Background(), "load config", "config", config.Config)

	r.Use(gin.Recovery())
	r.Use(ginmiddleware.CorsHandler())
	r.Use(ginmiddleware.Sentinel())

	if config.Config.Env.Profiles != "prod" {
		if config.Config.Env.Profiles == "test" {
			docs.SwaggerInfoapp.Host = "https://test-app-api.xiaoxixijz.com"
		}
		//配置swagger
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, func(c *ginSwagger.Config) {
			c.InstanceName = "app"
		}))
	}
	//健康检查
	r.GET("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	base := r.Group("/api/v1")
	base.Use(ginmiddleware.Logger())
	base.Use(ginmiddleware.RequestId(constant.RequestId))
	base.Use(ginmiddleware.WithContextDb(db))

	uploadApi := app.FileUpload{}
	uploadGroup := base.Group("/file")
	{
		uploadGroup.POST("/upload_image", uploadApi.UploadImage)
		uploadGroup.POST("/upload_base64_image", uploadApi.UploadBase64Image)
	}

	return r
}
