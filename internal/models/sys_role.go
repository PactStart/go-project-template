package models

type SysRole struct {
	SnowID
	Name        string `gorm:"column:name;type:varchar(20);comment:角色名称;unique;not null" json:"name"`
	Description string `gorm:"column:description;type:varchar(50);comment:角色描述;" json:"description"`
	BaseModel
}

type SysRoleUser struct {
	SnowID
	RoleID int64 `gorm:"column:role_id;uniqueIndex:uni_role_id_user_id;comment:角色ID" json:"roleId,omitempty"`
	UserID int64 `gorm:"column:user_id;uniqueIndex:uni_role_id_user_id;comment:用户ID" json:"userId,omitempty"`
	BaseModel
}

type SysRolePermission struct {
	SnowID
	RoleID       int64 `gorm:"column:role_id;uniqueIndex:uni_role_id_permission_id;comment:角色ID" json:"roleId,omitempty"`
	PermissionID int64 `gorm:"column:permission_id;uniqueIndex:uni_role_id_permission_id;comment:权限ID" json:"permissionId,omitempty"`
	BaseModel
}
