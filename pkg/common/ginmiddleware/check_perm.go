package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"orderin-server/internal/components"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/application"
	"orderin-server/pkg/common/constant"
	mcontext "orderin-server/pkg/common/context"
	"orderin-server/pkg/common/errs"
	"orderin-server/pkg/common/utils"
)

func CheckPermission(anonUrls []string, personalUrls []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsContain(c.Request.URL.Path, anonUrls) {
			c.Next()
		} else if mcontext.GetSuperAdmin(c) {
			c.Next()
		} else {
			opUserID, exists := c.Get(constant.OpUserID)
			if !exists {
				c.Abort()
			}
			uri := c.Request.URL.Path
			if utils.IsContain(uri, personalUrls) {
				c.Next()
			} else {
				exists = application.AppContext.GetComponent(components.COMPONENT_MY_CACHE).(*components.MyCache).HasPerm(opUserID.(int64), uri)
				if exists {
					c.Next()
				} else {
					api.GinError(c, errs.ErrNoPermission.Wrap())
					c.Abort()
				}
			}
		}
	}
}
