package dto

import (
	"orderin-server/internal/models"
	"orderin-server/pkg/common/dto"
)

type SysPermissionPageQueryReq struct {
	dto.Pagination `search:"-"`
	Type           string `json:"type" search:"type:exact;column:type;table:sys_permissions"`
	Anon           bool   `json:"anon" search:"type:exact;column:anon;table:sys_permissions"`
	Auth           bool   `json:"auth" search:"type:exact;column:auth;table:sys_permissions"`
	Keyword        string `json:"keyword" search:"type:contains;column:name,description;table:sys_permissions"`

	PermissionOrder
}

type PermissionOrder struct {
	PermissionIdOrder string `json:"idOrder" search:"type:order;column:id;table:sys_permissions" `
}

type SysPermissionAddReq struct {
	ParentID    int64  `json:"parentId,string" vd:"@:$>=0"`
	Type        string `json:"type" vd:"@:in($,'API','Page','Button')"`
	Anon        bool   `json:"anon"`
	Auth        bool   `json:"auth"`
	Name        string `json:"name" vd:"@:len($)>0 && len($)<100"`
	Description string `json:"description" vd:"@:len($)>0 && len($)<50"`
}

type SysPermissionUpdateReq struct {
	ID int64 `json:"id,string" vd:"@:$>0"`
	SysPermissionAddReq
}

type SysPermissionBatchAddReq struct {
	List []SysPermissionAddReq `json:"list" vd:"@:len($)>0"`
}

type SysPermissionDeleteReq struct {
	ID int64 `json:"id,string" vd:"@:$>0"`
}

type SysPermissionImportReq struct {
	SysPermissionAddReq
	Children []SysPermissionImportReq `json:"children"`
}

type SysPermissionTree struct {
	models.SysPermission
	Key      int64                `json:"key,string"`
	Title    string               `json:"title"`
	Children []*SysPermissionTree `json:"children"`
}
