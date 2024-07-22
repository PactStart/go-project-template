package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/internal/services"
	"orderin-server/pkg/common/api"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
)

type SysPermission struct {
	api.Api
}

// @Summary 添加权限
// @Description 添加权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param param body dto.SysPermissionAddReq false "权限信息"
// @Success 200 {object} api.Response
// @Router /sys/permission/add [post]
// @Security RequireLogin
func (e SysPermission) Add(context *gin.Context) {
	s := services.SysPermission{}
	req := dto.SysPermissionAddReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	permission := models.SysPermission{}
	utils.CopyStructFields(&permission, req)
	currentUserId := mcontext.GetOpUserID(context)
	permission.CreatedBy = &currentUserId

	_, err = s.Save(&permission)
	if err != nil {
		log.ZError(e.Context, "添加权限失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 删除权限
// @Description 根据id删除权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param param body dto.SysPermissionDeleteReq false "要删除的权限ID"
// @Success 200 {object} api.Response
// @Router /sys/permission/delete [post]
// @Security RequireLogin
func (e SysPermission) Delete(context *gin.Context) {
	s := services.SysPermission{}
	req := dto.SysPermissionDeleteReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	s.Delete(req.ID)
	e.OK(nil)
}

// @Summary 修改权限
// @Description 根据id修改权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param param body dto.SysPermissionUpdateReq false "权限信息"
// @Success 200 {object} api.Response
// @Router /sys/permission/update [post]
// @Security RequireLogin
func (e SysPermission) Update(context *gin.Context) {
	s := services.SysPermission{}
	req := dto.SysPermissionUpdateReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	permission := models.SysPermission{}
	utils.CopyStructFields(&permission, req)
	currentUserId := mcontext.GetOpUserID(context)
	permission.UpdatedBy = &currentUserId

	err = s.UpdateSelectiveById(permission.ID, permission)
	if err != nil {
		log.ZError(e.Context, "更新权限失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 分页查询权限
// @Description 分页查询权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param param body dto.SysPermissionPageQueryReq false "权限筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysPermission}}
// @Router /sys/permission/page_query [post]
// @Security RequireLogin
func (e SysPermission) PageQuery(context *gin.Context) {
	s := services.SysPermission{}
	req := dto.SysPermissionPageQueryReq{}
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
	list := make([]models.SysPermission, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 批量添加权限
// @Description 批量添加权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param param body dto.SysPermissionBatchAddReq false "权限筛选条件"
// @Success 200 {object} api.Response
// @Router /sys/permission/batch_add [post]
// @Security RequireLogin
func (e SysPermission) BatchAdd(context *gin.Context) {
	s := services.SysPermission{}
	req := dto.SysPermissionBatchAddReq{}
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
	permissions := utils.Batch[dto.SysPermissionAddReq, models.SysPermission](func(req dto.SysPermissionAddReq) models.SysPermission {
		permission := models.SysPermission{}
		utils.CopyStructFields(&permission, req)
		return permission
	}, req.List)
	err = s.BatchInsert(permissions)
	if err != nil {
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 批量导入权限
// @Description 批量导入权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param param body dto.SysPermissionImportReq false "权限树"
// @Success 200 {object} api.Response
// @Router /sys/permission/batch_import [post]
// @Security RequireLogin
func (e SysPermission) BatchImport(context *gin.Context) {
	s := services.SysPermission{}
	req := dto.SysPermissionImportReq{}
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
	err = s.Import(req)
	if err != nil {
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 获取权限树
// @Description 以树形结构呈现所有权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=[]dto.SysPermissionTree}
// @Router /sys/permission/bind_users [post]
// @Security RequireLogin
func (e SysPermission) GetTree(context *gin.Context) {
	s := services.SysPermission{}
	err := e.MakeContext(context).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	tree, err := s.GetTree()
	if err != nil {
		e.Error(err)
		return
	}
	e.OK(tree)
}
