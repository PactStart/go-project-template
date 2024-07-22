package dto

import (
	"orderin-server/pkg/common/customtypes"
	"orderin-server/pkg/common/dto"
)

type SysSmsSendReq struct {
	Phone string `json:"phone" vd:"@:phone($)"`
	Scene string `json:"scene" vd:"@:len($)>0"`
}

type SysSmsSendToMyselfReq struct {
	Scene string `json:"scene" vd:"@:len($)>0"`
}

type SysSmsLogPageQueryReq struct {
	dto.Pagination `search:"-"`
	Phone          string           `json:"phone" search:"type:exact;column:phone;table:sys_sms_logs"`
	TemplateCode   string           `json:"templateId" search:"type:exact;column:template_code;table:sys_sms_logs"`
	CreatedAtStart customtypes.Time `json:"createdAtStart" search:"type:gt;column:created_at;table:sys_sms_logs"`
	CreatedAtEnd   customtypes.Time `json:"createdAtEnd" search:"type:lt;column:created_at;table:sys_sms_logs"`

	SmsLogOrder
}

type SmsLogOrder struct {
	SMSLogIdOrder string `json:"idOrder" search:"type:order;column:id;table:sys_sms_logs" `
}
