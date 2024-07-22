package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/redis/go-redis/v9"
	"orderin-server/internal/components"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/application"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
)

type SysConfig struct {
	api.Api
}

// @Summary 添加配置
// @Description 添加配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param param body dto.SysConfigAddReq false "配置信息"
// @Success 200 {object} api.Response
// @Router /sys/config/add [post]
// @Security RequireLogin
func (e SysConfig) Add(context *gin.Context) {
	s := services.SysConfig{}
	req := dto.SysConfigAddReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	model := models.SysConfig{}
	utils.CopyStructFields(&model, req)
	currentUserId := mcontext.GetOpUserID(context)
	model.CreatedBy = &currentUserId

	err = s.Insert(&model)
	if err != nil {
		log.ZError(e.Context, "添加配置失败", err)
		e.Error(err)
		return
	}
	//发送消息
	components.Publish(application.AppContext.GetComponent(components.COMPONENT_REDIS).(redis.UniversalClient), components.CHANNEL_CONFIG_REFRESH, "{}")
	e.OK(nil)
}

// @Summary 修改配置
// @Description 根据id修改配置信息
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param param body dto.SysConfigUpdateReq false "配置信息"
// @Success 200 {object} api.Response
// @Router /sys/config/update [post]
// @Security RequireLogin
func (e SysConfig) Update(context *gin.Context) {
	s := services.SysConfig{}
	req := dto.SysConfigUpdateReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	model := models.SysConfig{}
	utils.CopyStructFields(&model, req)
	currentUserId := mcontext.GetOpUserID(context)
	model.UpdatedBy = &currentUserId

	err = s.UpdateSelectiveById(model.ID, model)
	if err != nil {
		log.ZError(e.Context, "修改配置失败", err)
		e.Error(err)
		return
	}
	//发送消息
	components.Publish(application.AppContext.GetComponent(components.COMPONENT_REDIS).(redis.UniversalClient), components.CHANNEL_CONFIG_REFRESH, "{}")
	e.OK(nil)
}

// @Summary 删除配置
// @Description 根据id删除配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param param body dto.SysConfigDeleteReq false "要删除的配置ID"
// @Success 200 {object} api.Response
// @Router /sys/config/delete [post]
// @Security RequireLogin
func (e SysConfig) Delete(context *gin.Context) {
	s := services.SysConfig{}
	req := dto.SysConfigDeleteReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	s.Delete(req.ID)
	//发送消息
	components.Publish(application.AppContext.GetComponent(components.COMPONENT_REDIS).(redis.UniversalClient), components.CHANNEL_CONFIG_REFRESH, "{}")
	e.OK(nil)
}

// @Summary 分页查询配置
// @Description 分页查询配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param param body dto.SysConfigPageQueryReq false "配置筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysConfig}}
// @Router /sys/config/page_query [post]
// @Security RequireLogin
func (e SysConfig) PageQuery(context *gin.Context) {
	s := services.SysConfig{}
	req := dto.SysConfigPageQueryReq{}
	err := e.MakeContext(context).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	list := make([]models.SysConfig, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 分页查询配置日志
// @Description 分页查询配置日志
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param param body dto.SysConfigLogPageQueryReq false "配置日志筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysConfigLog}}
// @Router /sys/config/page_query_log [post]
// @Security RequireLogin
func (e SysConfig) PageQueryLog(context *gin.Context) {
	s := services.SysConfig{}
	req := dto.SysConfigLogPageQueryReq{}
	err := e.MakeContext(context).
		MakeOrm().
		Bind(&req, binding.JSON).
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	list := make([]models.SysConfigLog, 0)
	var count int64

	err = s.PageQueryLog(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}
