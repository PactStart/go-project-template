package admin_api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net"
	"orderin-server/internal/components"
	"orderin-server/internal/models"
	"orderin-server/internal/routers"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/cache"
	"orderin-server/pkg/common/cmd"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
	"orderin-server/pkg/common/weixin"
	"os"
	"strconv"
	"strings"
)

var (
	AdminApiCmd = cmd.NewApiCmd(cmd.AdminApi)
)

func init() {
	AdminApiCmd.AddPortFlag()
	AdminApiCmd.AddPrometheusPortFlag()
	AdminApiCmd.AddApi(run)
}

func run(port int, proPort int) error {
	log.ZInfo(context.Background(), "admin-api port:", "port", port, "proPort", proPort)

	if port == 0 || proPort == 0 {
		err := "port or proPort is empty:" + strconv.Itoa(port) + "," + strconv.Itoa(proPort)
		log.ZError(context.Background(), err, nil)

		return fmt.Errorf(err)
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize Redis", err)
		return err
	}

	db, err := relation.NewGormDB()
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize mysql db", err)
		return err
	}
	err = initDatabase(db)
	if err != nil {
		log.ZError(context.Background(), "Failed to initialize database", err)
		return err
	}
	log.ZInfo(context.Background(), "db init success")

	router := routers.NewAdminGinRouter(db, rdb)
	log.ZInfo(context.Background(), "admin-api init routers success")

	weixin.OpenPlatformApp, err = weixin.NewOpenPlatformAppService()
	if err != nil || weixin.OpenPlatformApp == nil {
		log.ZError(context.Background(), "Failed to initialize OpenPlatformAppService", err)
		return err
	}
	components.RegisterComponents(db, rdb)

	var address string
	if config.Config.AdminApi.ListenIP != "" {
		address = net.JoinHostPort(config.Config.AdminApi.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	log.ZInfo(context.Background(), "start admin-api server", "address", address, "version", config.Version)

	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "admin-api run failed", err, "address", address)

		return err
	}

	return nil
}

func initDatabase(db *gorm.DB) error {
	//创建表
	if err := db.AutoMigrate(&models.SysUser{}, &models.SysRole{}, &models.SysPermission{}, &models.SysRoleUser{},
		&models.SysRolePermission{}, &models.SysConfig{}, &models.SysConfigLog{}, &models.SysDict{}, &models.SysDictItem{}, &models.SysSmsLog{},
		&models.WxAuthorizer{}, &models.WxAuthorizerMedia{}); err != nil {
		return err
	}
	//创建索引
	if !db.Migrator().HasIndex(&models.SysRoleUser{}, "uni_role_id_user_id") {
		db.Migrator().CreateIndex(&models.SysRoleUser{}, "uni_role_id_user_id")
	}
	if !db.Migrator().HasIndex(&models.SysRolePermission{}, "uni_role_id_permission_id") {
		db.Migrator().CreateIndex(&models.SysRolePermission{}, "uni_role_id_permission_id")
	}
	initSuperAdmin(db)
	err := initApiPermissions(db)
	if err != nil {
		log.ZError(context.Background(), "init api permissions occur error", err)
	}
	initDict(db)
	return nil
}

func initSuperAdmin(db *gorm.DB) {
	//创建默认的超级管理员用户
	s := services.SysUser{}
	s.Context = context.Background()
	s.Orm = db

	username := config.Config.SuperAdmin.Username
	password := config.Config.SuperAdmin.Password
	user, err := s.GetByUserName(username)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		log.ZInfo(context.Background(), "super admin isn't exist in db")
		user = &models.SysUser{
			Username:     username,
			RealName:     config.Config.SuperAdmin.RealName,
			Phone:        config.Config.SuperAdmin.Phone,
			PasswordSalt: utils.Md5(password),
			SuperAdmin:   true,
			Status:       1,
		}
		s.Insert(user)
		log.ZInfo(context.Background(), "success to create super admin", "username", username, "password", password)
	} else {
		user.PasswordSalt = utils.Md5(password)
		s.UpdateById(user.ID, *user)
	}
}

func initApiPermissions(db *gorm.DB) error {
	// 读取 Swagger 文档的 JSON 文件
	docBytes, err := os.ReadFile("./docs/admin-api/admin_swagger.json")
	if err != nil {
		log.ZError(context.Background(), "Read ./docs/admin-api/admin_swagger.json fail", err)
		return err
	}
	log.ZInfo(context.Background(), "Read ./docs/admin-api/admin_swagger.json success", "len", len(docBytes))

	// 解析 Swagger 文档的 JSON
	doc, err := loads.Analyzed(json.RawMessage(docBytes), "")
	if err != nil {
		log.ZError(context.Background(), "Parse ./admin-api/admin_swagger.json fail", err)
		return err
	}

	// 获取 Swagger 文档的根节点
	root := doc.Spec()

	//基础路径
	basePath := root.SwaggerProps.BasePath

	s := services.SysPermission{}
	s.Context = context.Background()
	s.Orm = db

	nameIdMap := make(map[string]int64)
	var ok = false
	var rootId int64 = 0
	var parentId int64 = 0

	rootId, err = s.Save(&models.SysPermission{
		Type:        "API",
		Anon:        false,
		Auth:        true,
		Name:        basePath + "/*",
		Description: "API权限",
	})
	if err != nil {
		log.ZError(context.Background(), "Write API permissions to DB fail", err)
		return err
	}
	// 遍历每个路径
	for path, pathItem := range root.Paths.Paths {
		var props spec.OperationProps
		if pathItem.PathItemProps.Post != nil {
			props = pathItem.PathItemProps.Post.OperationProps
		} else if pathItem.PathItemProps.Get != nil {
			props = pathItem.PathItemProps.Get.OperationProps
		} else {
			continue
		}
		fullPath := basePath + path

		parentId = 0
		if len(props.Tags) > 0 {
			lastIndex := strings.LastIndex(fullPath, "/")
			parentPath := fullPath[:lastIndex] + "/*"
			parentId, ok = nameIdMap[parentPath]
			if !ok {
				parentId, err = s.Save(&models.SysPermission{
					ParentID:    rootId,
					Type:        "API",
					Anon:        false,
					Auth:        true,
					Name:        parentPath,
					Description: props.Tags[0],
				})
				if err != nil {
					log.ZError(context.Background(), "Write API permissions to DB fail", err)
					return err
				}
				nameIdMap[parentPath] = parentId
			}
		}

		// 将接口信息写入数据库
		_, err2 := s.Save(&models.SysPermission{
			ParentID:    parentId,
			Type:        "API",
			Anon:        utils.IsContain(fullPath, routers.AnonUrls),
			Auth:        utils.IsContain(fullPath, routers.PersonalUrls),
			Name:        fullPath,
			Description: props.Summary,
		})
		if err2 != nil {
			log.ZError(context.Background(), "Write API permissions to DB fail", err)
			return err2
		}
	}

	return nil
}

func initDict(db *gorm.DB) {
	s := &services.SysDict{}
	s.Context = context.Background()
	s.Orm = db
	initSysDict(s)
}

func initSysDict(s *services.SysDict) {
	s.BatchAddDictAndItems("sys_user_status", "系统用户状态", []string{"正常", "禁用"}, []string{"1", "2"})
	s.BatchAddDictAndItems("sys_config_value_type", "配置数据类型", []string{"json", "bool", "number", "string"}, []string{"1", "2", "3", "4"})
	s.BatchAddDictAndItems("sys_sms_log_status", "短信发送状态", []string{"失败", "成功"}, []string{"1", "2"})
	s.BatchAddDictAndItems("sys_sms_log_template_code", "短信模板代码", []string{"验证码"}, []string{"SMS_296350564"})
}
