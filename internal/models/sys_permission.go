package models

type SysPermission struct {
	SnowID
	ParentID    int64  `gorm:"column:parent_id;comment:父权限ID" json:"parentId"`
	Type        string `gorm:"column:type;type:varchar(10);comment:类型，API,Page,Button" json:"type,omitempty"`
	Anon        bool   `gorm:"column:anon;type:tinyint(1);default:0;comment:可匿名访问;not null" json:"anon,omitempty"`
	Auth        bool   `gorm:"column:auth;type:tinyint(1);default:1;comment:是否需要鉴权;not null" json:"auth,omitempty"`
	Name        string `gorm:"column:name;type:varchar(100);comment:名称;unique;not null" json:"name,omitempty"`
	Description string `gorm:"column:description;type:varchar(50);comment:描述;" json:"description,omitempty"`
	BaseModel
}
