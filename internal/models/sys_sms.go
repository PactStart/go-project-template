package models

type SysSmsLog struct {
	SnowID
	TemplateCode string `gorm:"column:template_code;type:varchar(50);comment:短信模板;" json:"templateCode"`
	Vendor       string `gorm:"column:vendor;type:varchar(50);comment:短信服务商;" json:"vendor"`
	Scene        string `gorm:"column:scene;type:varchar(20);comment:短信场景;not null" json:"scene"`
	Phone        string `gorm:"column:phone;type:varchar(11);comment:手机号;not null" json:"phone"`
	Content      string `gorm:"column:content;type:varchar(255);comment:短信内容;not null" json:"content"`
	Status       string `gorm:"column:status;type:char(1);comment:状态;not null" json:"status"`
	Remark       string `gorm:"column:remark;type:varchar(100);comment:备注;not null" json:"remark"`
	CreatedModel
}
