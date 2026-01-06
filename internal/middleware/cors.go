package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			// 设置允许跨域的源
			c.Header("Access-Control-Allow-Origin", origin)

			// 设置允许的 HTTP 方法
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

			// 设置允许的请求头
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Origin, Accept")

			// 设置允许客户端携带凭证（如 cookies）
			c.Header("Access-Control-Allow-Credentials", "true")

			// 设置预检请求的缓存时间（单位为秒）
			c.Header("Access-Control-Max-Age", "86400")
		}

		// 如果是 OPTIONS 请求，直接返回 200 OK
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// 继续处理请求
		c.Next()
	}
}
