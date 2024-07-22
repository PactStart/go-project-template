package service

import (
	"context"
	"gorm.io/gorm"
)

type Service struct {
	Orm     *gorm.DB
	Context context.Context
}
