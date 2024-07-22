package relation

import (
	"context"
	"gorm.io/gorm"
)

type Tx interface {
	Transaction(fn func(tx any) error) error
}

type CtxTx interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

func NewTx(db *gorm.DB) Tx {
	return &_Gorm{tx: db}
}

type _Gorm struct {
	tx *gorm.DB
}

func (g *_Gorm) Transaction(fn func(tx any) error) error {
	return g.tx.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
