package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"orderin-server/pkg/common/log"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 计算处理时间
		cost := time.Since(start)
		// 打印请求信息
		log.ZInfo(c, "http request log", "uri", c.Request.RequestURI, "cost", cost.Milliseconds(), "method", c.Request.Method)
	}
}
