package models

type SysConfig struct {
	SnowID
	Name        string `gorm:"column:name;type:varchar(50);comment:配置名称;not null" json:"name"`
	Value       string `gorm:"column:value;type:tinytext;comment:配置值;" json:"value"`
	ValueType   string `gorm:"column:value_type;type:varchar(20);comment:值类型，number｜bool｜string｜json;not null" json:"valueType"`
	Description string `gorm:"column:description;type:varchar(100);comment:配置描述;not null" json:"description"`
	BaseModel
}

type SysConfigLog struct {
	SnowID
	ConfigID int64  `gorm:"column:config_id;comment:配置id;not null" json:"configId"`
	OldValue string `gorm:"column:old_value;type:tinytext;comment:旧值;" json:"oldValue"`
	NewValue string `gorm:"column:new_value;type:tinytext;comment:新值;" json:"newValue"`
	CreatedModel
}
