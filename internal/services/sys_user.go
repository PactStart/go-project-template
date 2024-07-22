package services

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
)

type SysUser struct {
	service.Service
}

func (e *SysUser) PageQuery(req *dto.SysUserPageQueryReq, list *[]models.SysUser, count *int64) error {
	tx := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		)
	if req.ExcludeRoleId > 0 {
		tx.Where("not exists (select 1 from sys_role_users where sys_users.id = sys_role_users.user_id and sys_role_users.role_id = ? )", req.ExcludeRoleId)
	}
	err := tx.Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_user fail", err)
		return err
	}
	return nil
}

func (e *SysUser) GetById(id int64) (*models.SysUser, error) {
	var user models.SysUser
	result := e.Orm.First(&user, id)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.AccountNotExistError, "用户不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, err
}

func (e *SysUser) Insert(user *models.SysUser) error {
	var err error
	var i int64
	err = e.Orm.Model(&user).Where("username = ?", user.Username).Count(&i).Error
	if err != nil {
		log.ZError(e.Context, "db error", err)
		return err
	}
	if i > 0 {
		log.ZError(e.Context, "username has been exist", err)
		err := errs.NewCodeError(errs.DuplicateKeyError, "用户名已存在")
		return err
	}
	err = e.Orm.Create(&user).Error
	return err
}

func (e *SysUser) GetByUserName(username string) (*models.SysUser, error) {
	var user models.SysUser
	err := e.Orm.Model(&user).Where("username = ?", username).First(&user).Error
	return &user, err
}

func (e *SysUser) GetByPhone(phone string) (*models.SysUser, error) {
	var user models.SysUser
	result := e.Orm.Model(&user).Where("phone = ?", phone).First(&user)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.AccountNotExistError, "手机号不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, err
}

func (e *SysUser) IsPhoneExist(phone string) (bool, error) {
	var err error
	var i int64
	err = e.Orm.Model(&models.SysUser{}).Where("phone = ?", phone).Count(&i).Error
	return i > 0, err
}

func (e *SysUser) GetByEmail(email string) (*models.SysUser, error) {
	var user models.SysUser
	result := e.Orm.Model(&user).Where("email = ?", email).First(&user)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.AccountNotExistError, "邮箱不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, err
}

func (e *SysUser) IsEmailExist(email string) (bool, error) {
	var err error
	var i int64
	err = e.Orm.Model(&models.SysUser{}).Where("email = ?", email).Count(&i).Error
	return i > 0, err
}

func (e *SysUser) GetByOpenID(openID string) (*models.SysUser, error) {
	var user models.SysUser
	result := e.Orm.Model(&user).Where("open_id = ?", openID).First(&user)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.AccountNotExistError, "openId不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &user, err
}

// 注意 gorm 默认情况下只会更新非零值字段
func (e *SysUser) UpdateById(id int64, user models.SysUser) error {
	user.ID = id
	return e.Orm.Model(&user).Updates(&user).Error
}

// 更新选定字段，即使零值也会更新，hook也会执行
func (e *SysUser) UpdateColumnsById(id int64, user models.SysUser, columns ...string) error {
	user.ID = id
	return e.Orm.Model(&user).Select(columns).Updates(user).Error
}

func (e *SysUser) GetRolesByUserId(id int64) ([]string, error) {
	var roleNames []string
	if err := e.Orm.Table("sys_roles").
		Joins("JOIN sys_role_users ON sys_roles.id = sys_role_users.role_id").
		Where("sys_role_users.user_id = ?", id).
		Pluck("sys_roles.name", &roleNames).Error; err != nil {
		return nil, err
	}
	return roleNames, nil
}

func (e *SysUser) GetPermissionsByUserId(id int64) ([]string, error) {
	var permissionNames []string
	if err := e.Orm.Table("sys_permissions").
		Joins("JOIN sys_role_permissions ON sys_permissions.id = sys_role_permissions.permission_id").
		Joins("JOIN sys_role_users ON sys_role_permissions.role_id = sys_role_users.role_id").
		Where("sys_role_users.user_id = ?", id).
		Pluck("sys_permissions.name", &permissionNames).Error; err != nil {
		return nil, err
	}
	return permissionNames, nil
}
