package models

type SysDict struct {
	IncrID
	Label       string `gorm:"column:label;type:varchar(50);comment:标签;" json:"label"`
	Name        string `gorm:"column:name;type:varchar(50);comment:字典名称;unique" json:"name"`
	Description string `gorm:"column:description;type:varchar(100);comment:字典描述;" json:"description"`
	BaseModel
	LogicalDeleted
}

type SysDictItem struct {
	IncrID
	DictName    string `gorm:"column:dict_name;type:varchar(50);comment:字典名称;" json:"dictName"`
	ItemLabel   string `gorm:"column:item_label;type:varchar(50);comment:标签;" json:"itemLabel"`
	ItemValue   string `gorm:"column:item_value;type:varchar(200);comment:值;" json:"itemValue"`
	Description string `gorm:"column:description;type:varchar(50);default:'';comment:字典项描述;" json:"description"`
	Status      int    `gorm:"column:status;type:int(11);comment:状态;" json:"status"`
	IsDefault   bool   `gorm:"column:is_default;type:tinyint(1);default:0;comment:是否默认;" json:"isDefault"`
	Sort        int    `gorm:"column:sort;type:int(11);comment:排序;" json:"sort"`
	BaseModel
}
