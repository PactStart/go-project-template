package services

import (
	"orderin-server/internal/dto"
	"orderin-server/internal/models"
	"orderin-server/pkg/common/db/relation"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/service"
)

type SysSms struct {
	service.Service
}

func (e SysSms) PageQuery(req *dto.SysSmsLogPageQueryReq, list *[]models.SysSmsLog, count *int64) error {
	err := e.Orm.
		Scopes(
			relation.MakeCondition(*req),
			relation.Paginate(req.GetPageSize(), req.GetPageIndex()),
		).Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		log.ZError(e.Context, "page query sys_sms_logs fail", err)
		return err
	}
	return nil
}

func (e SysSms) Insert(smsLog models.SysSmsLog) error {
	return e.Orm.Create(&smsLog).Error
}
