package relation

import (
	"context"

	"gorm.io/gorm"
)

type MetaDB struct {
	DB    *gorm.DB
	table any
}

func NewMetaDB(db *gorm.DB, table any) *MetaDB {
	return &MetaDB{
		DB:    db,
		table: table,
	}
}

func (g *MetaDB) GormDB(ctx context.Context) *gorm.DB {
	db := g.DB.WithContext(ctx).Model(g.table)
	return db
}
