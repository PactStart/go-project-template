package models

import (
	"gorm.io/gorm"
	"orderin-server/pkg/common/customtypes"
	"orderin-server/pkg/common/utils"
)

type SnowID struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:false" json:"id,string"`
}

type IncrID struct {
	ID int64 `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
}

type CreatedModel struct {
	CreatedAt customtypes.Time `gorm:"column:created_at;comment:创建时间;autoCreateTime;not null" json:"createdAt"`
	CreatedBy *int64           `gorm:"column:created_by;comment:创建人;" json:"createdBy,string"`
}

type BaseModel struct {
	CreatedAt customtypes.Time `gorm:"column:created_at;comment:创建时间;autoCreateTime;not null" json:"createdAt"`
	CreatedBy *int64           `gorm:"column:created_by;comment:创建人;" json:"createdBy,string"`
	UpdatedAt customtypes.Time `gorm:"column:updated_at;comment:更新时间;autoUpdateTime" json:"updatedAt"`
	UpdatedBy *int64           `gorm:"column:updated_by;comment:更新人" json:"updatedBy,string"`
}

type LogicalDeleted struct {
	Deleted   bool              `gorm:"column:deleted;type:tinyint(1);default:0;comment:是否删除;not null" json:"del,omitempty"`
	DeletedBy *int64            `gorm:"column:deleted_by;comment:删除人;" json:"deletedBy,string"`
	DeletedAt *customtypes.Time `gorm:"column:deleted_at;comment:删除时间;" json:"deletedAt"`
}

// 在模型的 BeforeCreate 钩子函数中生成雪花 ID
func (m *SnowID) BeforeCreate(tx *gorm.DB) error {
	m.ID = utils.GenID()
	return nil
}
