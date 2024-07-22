package dto

import (
	"orderin-server/pkg/common/customtypes"
	"orderin-server/pkg/common/dto"
)

type SysRoleAddReq struct {
	Name        string `json:"name" vd:"@:len($)>1 && len($)<20"`
	Description string `json:"description" vd:"@:len($)<50"`
}

type SysRoleUpdateReq struct {
	ID          int64  `json:"id,string" vd:"@:$>0"`
	Description string `json:"description" vd:"@:len($)<50"`
}

type SysRoleDeleteReq struct {
	ID int64 `json:"id,string" vd:"@:$>0"`
}

type SysRoleGetByNameReq struct {
	Name string `json:"name" vd:"@:len($)>1 && len($)<20"`
}

type SysRolePageQueryReq struct {
	dto.Pagination `search:"-"`
	Name           string `json:"name" search:"type:contains;column:name;table:sys_roles"`
	RoleOrder
}

type RoleOrder struct {
	RoleIdOrder string `json:"idOrder" search:"type:order;column:id;table:sys_roles" `
}

type SysRoleBindPermissionsReq struct {
	RoleId        int64                  `json:"roleId,string" vd:"@:$>0"`
	PermissionIds customtypes.Int64Slice `json:"permissionIds" vd:"@:len($)>0"`
}

type SysRoleGetPermissionTreeReq struct {
	RoleId int64 `json:"roleId,string" vd:"@:$>0"`
}

type SysRoleBindUsersReq struct {
	RoleId  int64                  `json:"roleId,string" vd:"@:$>0"`
	UserIds customtypes.Int64Slice `json:"userIds" vd:"@:len($)>0"`
}

type SysRoleUnbindUsersReq struct {
	RoleId  int64                  `json:"roleId,string" vd:"@:$>0"`
	UserIds customtypes.Int64Slice `json:"userIds" vd:"@:len($)>0"`
}

type SysRolePermissionTreeResp struct {
	OwnPermissionIds customtypes.Int64Slice `json:"ownPermissionIds"`
	PermissionTree   []*SysPermissionTree   `json:"permissionTree"`
}
