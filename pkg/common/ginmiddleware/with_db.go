package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithContextDb(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db.WithContext(c))
		c.Next()
	}
}
