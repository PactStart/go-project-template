package ginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"orderin-server/pkg/common/api"
	"orderin-server/pkg/common/constant"
	"orderin-server/pkg/common/errs"
)

// CorsHandler ginmiddleware cross-domain configuration.
func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar",
		) // Cross-domain key settings allow browsers to resolve.
		c.Header(
			"Access-Control-Max-Age",
			"172800",
		) // Cache request information in seconds.
		c.Header(
			"Access-Control-Allow-Credentials",
			"false",
		) //  Whether cross-domain requests need to carry cookie information, the default setting is true.
		c.Header(
			"content-type",
			"application/json",
		) // Set the return format to json.
		// Release all option pre-requests
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, "Options Request!")
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequiredHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在这里定义你要校验的请求头
		requiredHeaders := []string{
			constant.Token,
		}

		for _, header := range requiredHeaders {
			if value := c.GetHeader(header); value == "" {
				// 如果请求头缺失，则返回错误响应
				err := errors.New(fmt.Sprintf("Missing required header: %s", header))
				api.GinError(c, errs.ErrArgs.Wrap(err.Error()))
				c.Abort()
				return
			}
		}
		// 请求头校验通过，继续处理请求
		c.Next()
	}
}
