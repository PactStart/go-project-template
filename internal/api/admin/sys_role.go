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

type SysRole struct {
	api.Api
}

// @Summary 添加角色
// @Description 添加角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleAddReq false "角色信息"
// @Success 200 {object} api.Response
// @Router /sys/role/add [post]
// @Security RequireLogin
func (e SysRole) Add(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleAddReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}

	role := models.SysRole{}
	utils.CopyStructFields(&role, req)
	currentUserId := mcontext.GetOpUserID(context)
	role.CreatedBy = &currentUserId

	err = s.Insert(&role)
	if err != nil {
		log.ZError(e.Context, "添加角色失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 修改角色
// @Description 根据id修改角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleUpdateReq false "角色信息"
// @Success 200 {object} api.Response
// @Router /sys/role/update [post]
// @Security RequireLogin
func (e SysRole) Update(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleUpdateReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	role := models.SysRole{}
	utils.CopyStructFields(&role, req)
	currentUserId := mcontext.GetOpUserID(context)
	role.UpdatedBy = &currentUserId

	err = s.UpdateSelectiveById(role.ID, role)
	if err != nil {
		log.ZError(e.Context, "更新角色失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 删除角色
// @Description 根据id删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleDeleteReq false "要删除的角色ID"
// @Success 200 {object} api.Response
// @Router /sys/role/delete [post]
// @Security RequireLogin
func (e SysRole) Delete(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleDeleteReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	s.Delete(req.ID)
	e.OK(nil)
}

// @Summary 根据名字获取角色
// @Description 根据名字获取角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleGetByNameReq false "角色名称"
// @Success 200 {object} api.Response
// @Router /sys/role/get_by_name [post]
// @Security RequireLogin
func (e SysRole) GetByName(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleGetByNameReq{}
	err := e.MakeContext(context).Bind(&req, binding.JSON).MakeOrm().MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	role, err := s.GetByName(req.Name)
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	if role.ID > 0 {
		e.OK(role)
	} else {
		e.OK(nil)
	}
}

// @Summary 分页查询角色
// @Description 分页查询角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRolePageQueryReq false "角色筛选条件"
// @Success 200 {object} api.Response{data=api.PageData{List=models.SysRole}}
// @Router /sys/role/page_query [post]
// @Security RequireLogin
func (e SysRole) PageQuery(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRolePageQueryReq{}
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
	list := make([]models.SysRole, 0)
	var count int64

	err = s.PageQuery(&req, &list, &count)
	if err != nil {
		e.Error(err)
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize())
}

// @Summary 给角色授予权限
// @Description 批量关联权限给角色，会覆盖之前角色拥有的权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleBindPermissionsReq true "角色id和权限id数组"
// @Success 200 {object} api.Response
// @Router /sys/role/bind_permissions [post]
// @Security RequireLogin
func (e SysRole) BindPermissions(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleBindPermissionsReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	err = s.BindPermissions(req)
	if err != nil {
		log.ZError(e.Context, "角色绑定权限失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 获取角色拥有的权限树
// @Description 获取角色拥有的权限树
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleGetPermissionTreeReq true "角色id"
// @Success 200 {object} api.Response{data=dto.SysRolePermissionTreeResp}
// @Router /sys/role/get_permission_tree [post]
// @Security RequireLogin
func (e SysRole) GetPermissionTree(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleGetPermissionTreeReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	ownPermissionIds, err := s.GetPermissionIdsByRoleId(req.RoleId)
	if err != nil {
		log.ZError(e.Context, "获取角色拥有的权限失败", err)
		e.Error(err)
		return
	}
	permissionService := services.SysPermission{}
	permissionService.Context = context
	permissionService.Orm = e.Orm
	permissionTree, err := permissionService.GetTree()
	if err != nil {
		log.ZError(e.Context, "获取权限树失败", err)
		e.Error(err)
		return
	}

	resp := dto.SysRolePermissionTreeResp{
		OwnPermissionIds: ownPermissionIds,
		PermissionTree:   permissionTree,
	}
	e.OK(resp)
}

// @Summary 授权角色给用户
// @Description 角色关联用户，关联后用户拥有该角色和该角色的所有权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleBindPermissionsReq true "需要绑定的角色id和用户id数组"
// @Success 200 {object} api.Response
// @Router /sys/role/bind_users [post]
// @Security RequireLogin
func (e SysRole) BindUsers(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleBindUsersReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	err = s.BindUsers(req)
	if err != nil {
		log.ZError(e.Context, "角色绑定用户失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}

// @Summary 取消授权角色给用户
// @Description 取消角色关联用户，取消后用户不再拥有该角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param param body dto.SysRoleUnbindUsersReq true "需要解绑的角色id和用户id数组"
// @Success 200 {object} api.Response
// @Router /sys/role/unbind_users [post]
// @Security RequireLogin
func (e SysRole) UnbindUsers(context *gin.Context) {
	s := services.SysRole{}
	req := dto.SysRoleUnbindUsersReq{}
	err := e.MakeContext(context).MakeOrm().Bind(&req, binding.JSON).MakeService(&s.Service).Errors
	if err != nil {
		log.ZError(e.Context, "", err)
		e.Error(err)
		return
	}
	err = s.UnBindUsers(req)
	if err != nil {
		log.ZError(e.Context, "角色解绑用户失败", err)
		e.Error(err)
		return
	}
	e.OK(nil)
}
