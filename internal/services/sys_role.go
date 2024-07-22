package services

import (
	"errors"
	"gorm.io/gorm"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
	"orderin-server/pkg/common/utils"
)

type SysRole struct {
	service.Service
}

func (e SysRole) PageQuery(req *dto.SysRolePageQueryReq, list *[]models.SysRole, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_roles fail", err)
		return err
	}
	return nil
}

func (e SysRole) Insert(role *models.SysRole) error {
	var err error
	var i int64
	err = e.Orm.Model(role).Where("name = ?", role.Name).Count(&i).Error
	if err != nil {
		log.ZError(e.Context, "db error", err)
		return err
	}
	if i > 0 {
		log.ZError(e.Context, "role has been exist", err)
		err := errs.NewCodeError(errs.DuplicateKeyError, "角色已存在")
		return err
	}
	err = e.Orm.Create(&role).Error
	return err
}

func (e SysRole) UpdateSelectiveById(id int64, role models.SysRole) error {
	role.ID = id
	return e.Orm.Model(&role).Updates(&role).Error
}

func (e SysRole) Delete(id int64) {
	e.Orm.Delete(models.SysRole{}, id)
}

func (e SysRole) GetByName(name string) (*models.SysRole, error) {
	var role models.SysRole
	result := e.Orm.Model(&role).Where("name = ?", name).First(&role)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &role, err
}

func (e SysRole) BindUsers(req dto.SysRoleBindUsersReq) error {
	roleUsers := utils.Batch[int64, models.SysRoleUser](func(userId int64) models.SysRoleUser {
		roleUser := models.SysRoleUser{
			RoleID: req.RoleId,
			UserID: userId,
		}
		currentUserId := mcontext.GetOpUserID(e.Context)
		roleUser.CreatedBy = &currentUserId
		return roleUser
	}, req.UserIds)
	return e.Orm.CreateInBatches(roleUsers, len(roleUsers)).Error
}

func (e SysRole) UnBindUsers(req dto.SysRoleUnbindUsersReq) error {
	return e.Orm.Where("role_id = ? AND user_id IN (?)", req.RoleId, req.UserIds).Delete(&models.SysRoleUser{}).Error
}

func (e SysRole) GetPermissionIdsByRoleId(roleId int64) ([]int64, error) {
	var permissionIds []int64
	if err := e.Orm.Table("sys_permissions").
		Joins("JOIN sys_role_permissions ON sys_permissions.id = sys_role_permissions.permission_id").
		Where("sys_role_permissions.role_id = ?", roleId).Pluck("sys_permissions.id", &permissionIds).Error; err != nil {
		return nil, err
	}
	return permissionIds, nil
}

func (e SysRole) BindPermissions(req dto.SysRoleBindPermissionsReq) error {
	// 开始事务
	tx := e.Orm.Begin()

	// 使用 defer 设置事务的回滚或提交
	defer func() {
		if r := recover(); r != nil {
			// 发生了 panic，回滚事务
			tx.Rollback()
			log.ZError(e.Context, "Transaction rolled back", nil, "err", r)
		} else if tx.Error != nil {
			// 发生了错误，回滚事务
			tx.Rollback()
			log.ZError(e.Context, "Transaction rolled back", tx.Error)

		} else {
			// 没有发生错误，提交事务
			tx.Commit()
			log.ZInfo(e.Context, "Transaction committed")

		}
	}()
	var err error
	err = tx.Where("role_id = ?", req.RoleId).Delete(&models.SysRolePermission{}).Error
	if err != nil {
		log.ZError(e.Context, "删除当前角色已有的权限失败", err)
		return err
	}
	rolePermissions := utils.Batch[int64, models.SysRolePermission](func(permissionId int64) models.SysRolePermission {
		rolePermission := models.SysRolePermission{
			RoleID:       req.RoleId,
			PermissionID: permissionId,
		}
		currentUserId := mcontext.GetOpUserID(e.Context)
		rolePermission.CreatedBy = &currentUserId
		return rolePermission
	}, req.PermissionIds)
	err = tx.CreateInBatches(rolePermissions, len(rolePermissions)).Error
	if err != nil {
		log.ZError(e.Context, "给角色绑定权限失败", err)
		return err
	}
	return nil

}
