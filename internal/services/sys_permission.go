package services

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
	"orderin-server/pkg/common/utils"
)

type SysPermission struct {
	service.Service
}

func (e SysPermission) PageQuery(req *dto.SysPermissionPageQueryReq, list *[]models.SysPermission, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_permission fail", err)
		return err
	}
	return nil
}

func (e SysPermission) Save(permission *models.SysPermission) (int64, error) {
	tx := e.Orm.Begin()
	id, err := e.DoSave(tx, permission)
	if err == nil {
		tx.Commit()
	}
	return id, err
}

func (e SysPermission) DoSave(tx *gorm.DB, permission *models.SysPermission) (int64, error) {
	var err error
	var id int64 = 0
	var dbRecord models.SysPermission
	result := tx.Model(&permission).Where("name = ?", permission.Name).Find(&dbRecord)
	err = result.Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.ZError(e.Context, "db error", err)
			return 0, err
		}
	}
	if result.RowsAffected == 0 {
		err = tx.Create(&permission).Error
		if err != nil {
			return 0, err
		}
		id = permission.ID
	} else {
		dbRecord.Type = permission.Type
		dbRecord.ParentID = permission.ParentID
		dbRecord.Anon = permission.Anon
		dbRecord.Auth = permission.Auth
		dbRecord.Description = permission.Description
		dbRecord.UpdatedBy = permission.CreatedBy

		tx.Save(&dbRecord)
		id = dbRecord.ID
	}
	return id, err
}

func (e *SysPermission) UpdateSelectiveById(id int64, permission models.SysPermission) error {
	permission.ID = id
	return e.Orm.Model(&permission).Updates(&permission).Error
}

func (e SysPermission) Delete(id int64) {
	e.Orm.Delete(models.SysPermission{}, id)
}

func (e SysPermission) GetTree() ([]*dto.SysPermissionTree, error) {
	var list []models.SysPermission
	err := e.Orm.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return buildPermTree(list), nil
}

func (e SysPermission) BatchInsert(permissions []models.SysPermission) error {
	// 开始事务
	tx := e.Orm.Begin()

	// 使用 defer 设置事务的回滚或提交
	defer func() {
		if err := recover(); err != nil {
			// 发生了 panic，回滚事务
			tx.Rollback()
			log.ZError(e.Context, "Transaction rolled back", nil, "err", err)
		} else if tx.Error != nil {
			// 发生了错误，回滚事务
			tx.Rollback()
			log.ZError(e.Context, "Transaction rolled back", tx.Error)

		} else {
			// 没有发生错误，提交事务
			tx.Commit()
			log.ZInfo(e.Context, "Transaction committed", tx.Error)

		}
	}()

	var err error
	for _, permission := range permissions {
		_, err = e.Save(&permission)
		if err != nil {
			return err
		}
	}
	return err
}

func (e SysPermission) Import(req dto.SysPermissionImportReq) error {
	// 开始事务
	tx := e.Orm.Begin()

	// 使用 defer 设置事务的回滚或提交
	defer func() {
		if err := recover(); err != nil {
			// 发生了 panic，回滚事务
			tx.Rollback()
			log.ZError(e.Context, "Transaction rolled back", err.(error))
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

	err := e.WalkAndSave(tx, nil, req)
	if err != nil {
		return err
	}
	return nil
}

func (e SysPermission) WalkAndSave(tx *gorm.DB, parentId *int64, req dto.SysPermissionImportReq) error {
	permission := models.SysPermission{}
	utils.CopyStructFields(&permission, req)
	currentUserId := mcontext.GetOpUserID(e.Context)
	permission.CreatedBy = &currentUserId

	if parentId != nil {
		permission.ParentID = *parentId
	}

	selfId, err := e.DoSave(tx, &permission)

	if len(req.Children) > 0 {
		for _, child := range req.Children {
			e.WalkAndSave(tx, &selfId, child)
		}
	}
	return err
}

// 辅助函数：递归构建权限树
func buildPermTree(perms []models.SysPermission) []*dto.SysPermissionTree {
	permMap := make(map[int64]*dto.SysPermissionTree)

	for _, perm := range perms {
		permTree := dto.SysPermissionTree{
			SysPermission: perm,
			Key:           perm.ID,
			Title:         perm.Description + " (" + perm.Name + ")",
		}
		permMap[perm.ID] = &permTree
	}

	var tree []*dto.SysPermissionTree
	for _, perm := range perms {
		if parent, ok := permMap[perm.ParentID]; ok {
			parent.Children = append(parent.Children, permMap[perm.ID])
		} else {
			tree = append(tree, permMap[perm.ID])
		}
	}
	if tree == nil {
		tree = make([]*dto.SysPermissionTree, 0)
	}
	return tree
}
