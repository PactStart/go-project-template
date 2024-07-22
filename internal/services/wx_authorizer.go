package services

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/pkg/common/config"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
)

type WxAuthorizer struct {
	service.Service
}

func (e WxAuthorizer) Save(authorizer *models.WxAuthorizer) error {
	var dbRecord models.WxAuthorizer
	tx := e.Orm.Where("component_appid = ? and authorizer_appid = ?", authorizer.ComponentAppid, authorizer.AuthorizerAppid).First(&dbRecord)
	err := tx.Error
	if err != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return tx.Error
		} else {
			return e.Orm.Create(authorizer).Error
		}
	}
	if tx.RowsAffected == 0 {
		return e.Orm.Create(authorizer).Error
	} else {
		authorizer.ID = dbRecord.ID
		return e.Orm.Updates(authorizer).Error
	}
}

func (e WxAuthorizer) GetById(id int64) (*models.WxAuthorizer, error) {
	var authorizer models.WxAuthorizer
	result := e.Orm.First(&authorizer, id)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.RecordNotFoundError, "授权方不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &authorizer, err
}

func (e WxAuthorizer) GetByAppID(appID string) (*models.WxAuthorizer, error) {
	var authorizer models.WxAuthorizer
	result := e.Orm.Where("component_appid = ? and authorizer_appid = ?", config.Config.WxOpenPlatform.AppID, appID).First(&authorizer)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewCodeError(errs.RecordNotFoundError, "授权方不存在")
		} else {
			return nil, err
		}
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &authorizer, err
}

func (e WxAuthorizer) GetAll() (*[]models.WxAuthorizer, error) {
	authorizers := make([]models.WxAuthorizer, 0)
	err := e.Orm.Find(&authorizers).Error
	return &authorizers, err
}

func (e WxAuthorizer) PageQuery(req *dto.WxAuthorizerPageQueryReq, list *[]models.WxAuthorizer, count *int64) error {
	err := e.Orm.Scopes(relation.MakeCondition(*req),
		relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
	).Find(list).Limit(-1).Offset(-1).Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query wx_authorizers fail", err)
		return err
	}
	return nil
}

func (e WxAuthorizer) UpdateById(id int64, authorizer *models.WxAuthorizer) error {
	authorizer.ID = id
	return e.Orm.Model(&authorizer).Save(&authorizer).Error
}

func (e WxAuthorizer) UpdateDefaultKfAccountByAppId(account string, appId string) error {
	return e.Orm.Raw("update wx_authorizers set default_kf_account = ? where authorizer_appid = ?", account, appId).Error
}

func (e WxAuthorizer) GetMedia(appId string, category string) (*models.WxAuthorizerMedia, error) {
	var dbRecord models.WxAuthorizerMedia
	err := e.Orm.Where("app_id = ? and category = ?", appId, category).First(&dbRecord).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &dbRecord, nil
}

func (e WxAuthorizer) InsertMedia(media *models.WxAuthorizerMedia) error {
	return e.Orm.Create(media).Error
}
