package services

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
)

type SysConfig struct {
	service.Service
}

func (e SysConfig) Insert(config *models.SysConfig) error {
	var err error
	var i int64
	err = e.Orm.Model(&config).Where("name = ?", config.Name).Count(&i).Error
	if err != nil {
		log.ZError(e.Context, "db error", err)
		return err
	}
	if i > 0 {
		log.ZError(e.Context, "config has been exist", err)
		err := errs.NewCodeError(errs.DuplicateKeyError, "配置已存在")
		return err
	}
	err = e.Orm.Create(&config).Error
	return err
}

func (e SysConfig) UpdateSelectiveById(id int64, config models.SysConfig) error {
	var dbRecord models.SysConfig
	var err error

	result := e.Orm.First(&dbRecord, id)
	err = result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NewCodeError(errs.RecordNotFoundError, "配置不存在")
		}
		return err
	}
	if result.RowsAffected == 0 {
		return errs.NewCodeError(errs.RecordNotFoundError, "配置不存在")
	}

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

	config.ID = id
	err = tx.Model(&config).Updates(&config).Error
	if err != nil {
		return err
	}

	if dbRecord.Value != config.Value {
		log := models.SysConfigLog{
			ConfigID: id,
			OldValue: dbRecord.Value,
			NewValue: config.Value,
		}
		currentUserId := mcontext.GetOpUserID(e.Context)
		log.CreatedBy = &currentUserId
		err = tx.Create(&log).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (e SysConfig) Delete(id int64) {
	e.Orm.Delete(models.SysConfig{}, id)
}

func (e SysConfig) PageQuery(req *dto.SysConfigPageQueryReq, list *[]models.SysConfig, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_configs fail", err)
		return err
	}
	return nil
}

func (e SysConfig) PageQueryLog(req *dto.SysConfigLogPageQueryReq, list *[]models.SysConfigLog, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_config_logs fail", err)
		return err
	}
	return nil
}

func (e SysConfig) GetAll() (*[]models.SysConfig, error) {
	var configs []models.SysConfig
	err := e.Orm.Find(&configs).Error
	return &configs, err
}
