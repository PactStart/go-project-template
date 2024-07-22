package routers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"net/http"
	docs "orderin-server/docs/admin-api"
	"orderin-server/internal/api/admin"
	"orderin-server/pkg/common/auth"
	"orderin-server/pkg/common/cache"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/ginmiddleware"
	"orderin-server/pkg/common/log"
)

var (
	AnonUrls = []string{
		"/api/v1/sys/user/login_by_password",
		"/api/v1/sys/user/login_by_phone",
		"/api/v1/sys/user/bind_wechat",
		"/api/v1/sys/user/get_wechat_login_url",
		"/api/v1/sys/user/login_by_wechat",
		"/api/v1/sys/user/get_wechat_login_result",
		"/api/v1/sys/sms/send",
		"/api/v1/wx/authorizer/get_authorize_url",
		"/api/v1/wx/authorizer/authorize/redirect",
	}
	PersonalUrls = []string{
		"/api/v1/sys/user/logout",
		"/api/v1/sys/user/update_password",
		"/api/v1/sys/user/get_personal_info",
		"/api/v1/sys/user/perfect_info",
		"/api/v1/sys/user/bind_mfa",
		"/api/v1/sys/user/get_oauth2_url",
	}
)

// @title Admin API
// @version 1.0.0
// @description xxx管理后台API文档

// @contact.name xxx
// @contact.url http://xxx.com
// @contact.email xxx@qq.com

// @securityDefinitions.apikey RequireLogin
// @in header
// @name token

// @host 127.0.0.1:10000
// @BasePath /api/v1
func NewAdminGinRouter(db *gorm.DB, rdb redis.UniversalClient) *gin.Engine {
	if config.Config.Env.Profiles == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	//	_ = v.RegisterValidation("required_if", ginmiddleware.RequiredIf)
	//}
	log.ZInfo(context.Background(), "load config", "config", config.Config)

	authService := auth.NewAuthService(
		cache.NewTokenCache(rdb),
		config.Config.Secret,
		config.Config.TokenPolicy.Expire,
	)

	r.Use(gin.Recovery())
	r.Use(ginmiddleware.CorsHandler())
	r.Use(ginmiddleware.Sentinel())

	if config.Config.Env.Profiles != "prod" {
		if config.Config.Env.Profiles == "test" {
			docs.SwaggerInfoadmin.Host = "https://test-admin-api.xxx.com"
		}
		//配置swagger
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, func(c *ginSwagger.Config) {
			c.InstanceName = "admin"
		}))
	}
	//健康检查
	r.GET("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	//微信开放平台校验文件
	r.GET("/wI2oB2J91M.txt", func(c *gin.Context) {
		c.String(200, "baec944df17b60f36e2f6e519e9db70b")
	})

	base := r.Group("/api/v1")
	base.Use(ginmiddleware.Logger())
	base.Use(ginmiddleware.RequestId(constant.RequestId))
	base.Use(ginmiddleware.WithContextDb(db))
	//解析token，排除匿名接口
	base.Use(ginmiddleware.ParseToken(authService, rdb, AnonUrls))
	//验证权限，排除个人接口
	base.Use(ginmiddleware.CheckPermission(AnonUrls, PersonalUrls))

	verifyCodeCache := cache.NewVerifyCodeCache(rdb)
	userApi := admin.SysUser{
		AuthService:     authService,
		VerifyCodeCache: verifyCodeCache,
	}

	uploadApi := admin.FileUpload{}
	uploadGroup := base.Group("/file")
	{
		uploadGroup.POST("/upload_image", uploadApi.UploadImage)
		uploadGroup.POST("/upload_base64_image", uploadApi.UploadBase64Image)
		uploadGroup.POST("/upload_blob_image", uploadApi.UploadBlobImage)
		uploadGroup.POST("/download_and_store", uploadApi.DownloadAndStore)
	}

	userGroup := base.Group("/sys/user")
	{
		userGroup.POST("/login_by_password", userApi.LoginByPassword)
		userGroup.POST("/login_by_phone", userApi.LoginByPhone)
		userGroup.POST("/logout", userApi.Logout)

		userGroup.POST("/page_query", userApi.PageQuery)
		userGroup.POST("/get_by_id", userApi.GetById)
		userGroup.POST("/add", userApi.Add)
		userGroup.POST("/update_status", userApi.UpdateStatus)
		userGroup.POST("/reset_password", userApi.ResetPassword)

		userGroup.POST("/get_personal_info", userApi.GetPersonalInfo)
		userGroup.POST("/perfect_info", userApi.PerfectInfo)
		userGroup.POST("/bind_email", userApi.BindEmail)
		userGroup.POST("/bind_phone", userApi.BindPhone)
		userGroup.POST("/update_password", userApi.UpdatePassword)
		userGroup.POST("/change_password", userApi.ChangePassword)
		userGroup.POST("/generate_mfa_device", userApi.GenerateMFADevice)
		userGroup.POST("/bind_mfa_device", userApi.BindMFADevice)
		userGroup.POST("/unbind_mfa_device", userApi.UnBindMFADevice)
		userGroup.POST("/switch_enable_mfa", userApi.SwitchEnableMFA)

		userGroup.POST("/get_oauth2_url", userApi.GetOAuth2Url)
		userGroup.GET("/bind_wechat", userApi.BindWechat)
		userGroup.POST("/unbind_wechat", userApi.UnBindWechat)

		userGroup.POST("/get_wechat_login_url", userApi.GetWxLoginUrl)
		userGroup.GET("/login_by_wechat", userApi.LoginByWechat)
		userGroup.POST("/get_wechat_login_result", userApi.GetWechatLoginResult)

		userGroup.POST("/send_verify_email", userApi.SendVerifyEmail)

		userGroup.POST("/should_check_identity", userApi.ShouldCheckIdentity)
		userGroup.POST("/check_identity", userApi.CheckIdentity)
	}

	permissionApi := admin.SysPermission{}
	permissionGroup := base.Group("/sys/permission")
	{
		permissionGroup.POST("/add", permissionApi.Add)
		permissionGroup.POST("/delete", permissionApi.Delete)
		permissionGroup.POST("/update", permissionApi.Update)
		permissionGroup.POST("/page_query", permissionApi.PageQuery)
		permissionGroup.POST("/batch_add", permissionApi.BatchAdd)
		permissionGroup.POST("/batch_import", permissionApi.BatchImport)
		permissionGroup.POST("/get_tree", permissionApi.GetTree)
	}

	roleApi := admin.SysRole{}
	roleGroup := base.Group("/sys/role")
	{
		roleGroup.POST("/add", roleApi.Add)
		roleGroup.POST("/update", roleApi.Update)
		roleGroup.POST("/delete", roleApi.Delete)
		roleGroup.POST("/get_by_name", roleApi.GetByName)
		roleGroup.POST("/page_query", roleApi.PageQuery)
		roleGroup.POST("/bind_permissions", roleApi.BindPermissions)
		roleGroup.POST("/get_permission_tree", roleApi.GetPermissionTree)
		roleGroup.POST("/bind_users", roleApi.BindUsers)
		roleGroup.POST("/unbind_users", roleApi.UnbindUsers)
	}

	configApi := admin.SysConfig{}
	configGroup := base.Group("/sys/config")
	{
		configGroup.POST("/add", configApi.Add)
		configGroup.POST("/update", configApi.Update)
		configGroup.POST("/delete", configApi.Delete)
		configGroup.POST("/page_query", configApi.PageQuery)
		configGroup.POST("/page_query_log", configApi.PageQueryLog)
	}

	dictApi := admin.SysDict{}
	dictGroup := base.Group("/sys/dict")
	{
		dictGroup.POST("/add", dictApi.Add)
		dictGroup.POST("/update", dictApi.Update)
		dictGroup.POST("/delete", dictApi.Delete)
		dictGroup.POST("/page_query", dictApi.PageQuery)
		dictGroup.POST("/item/add", dictApi.AddItem)
		dictGroup.POST("/item/update", dictApi.UpdateItem)
		dictGroup.POST("/item/page_query", dictApi.PageQueryItem)
		dictGroup.POST("/item/get_by_name", dictApi.GetItemsByName)
		dictGroup.POST("/item/batch_get_by_names", dictApi.BatchGetItemsByNames)
	}

	smsApi := admin.SysSms{VerifyCodeCache: verifyCodeCache}
	smsGroup := base.Group("/sys/sms")
	{
		smsGroup.POST("/send", smsApi.Send)
		smsGroup.POST("/send_to_myself", smsApi.SendToMyself)
		smsGroup.POST("/page_query", smsApi.PageQuery)
	}

	openWxApi := admin.WxOpenPlatformServer{}
	{
		// auth callback
		base.POST("/callback", openWxApi.Callback)
		base.POST("/callback/:appID", openWxApi.CallbackWithApp)

	}

	authorizerApi := admin.WxAuthorizer{}
	authorizerGroup := base.Group("/wx/authorizer")
	{
		authorizerGroup.POST("/page_query", authorizerApi.PageQuery)
		authorizerGroup.POST("/get_authorize_url", authorizerApi.GetAuthorizeUrl)
		authorizerGroup.GET("/authorize/redirect", authorizerApi.AuthorizeRedirect)
		authorizerGroup.POST("/sync_info", authorizerApi.SyncInfo)

		authorizerGroup.POST("/create_kf_account", authorizerApi.CreateKfAccount)
		authorizerGroup.POST("/get_all_kf_accounts", authorizerApi.GetAllKfAccounts)
		authorizerGroup.POST("/delete_kf_account", authorizerApi.DeleteKfAccount)
		authorizerGroup.POST("/set_default_kf_account", authorizerApi.SetDefaultKfAccount)
		authorizerGroup.POST("/delete_menu", authorizerApi.DeleteMenu)
		authorizerGroup.POST("/get_menu", authorizerApi.GetMenu)

	}

	return r
}
