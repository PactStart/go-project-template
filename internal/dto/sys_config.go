package dto

import "orderin-server/pkg/common/dto"

type SysConfigAddReq struct {
	Name        string `json:"name" vd:"@:len($)>0 && len($)<50"`
	Value       string `json:"value" vd:"@:len($)>0"`
	ValueType   string `json:"ValueType" vd:"@:len($)>0"`
	Description string `json:"description" vd:"@:len($)<100"`
}

type SysConfigUpdateReq struct {
	ID          int64  `json:"id,string" vd:"@:$>0"`
	Value       string `json:"value" vd:"@:len($)>0"`
	Description string `json:"description" vd:"@:len($)<100"`
}

type SysConfigDeleteReq struct {
	ID int64 `json:"id,string" vd:"@:$>0"`
}

type SysConfigPageQueryReq struct {
	dto.Pagination `search:"-"`
	Keyword        string `json:"keyword" search:"type:contains;column:name,description;table:sys_configs"`
	ConfigOrder
}

type ConfigOrder struct {
	ConfigIdOrder string `json:"idOrder" search:"type:order;column:id;table:sys_configs" `
}

type SysConfigLogPageQueryReq struct {
	dto.Pagination `search:"-"`
	ConfigID       int64 `json:"configId,string" vd:"@:$>0" search:"type:exact;column:config_id;table:sys_config_logs"`
}
