package middleware

import (
	"chain/internal/utils"

	"github.com/gin-gonic/gin"
)

// Gin中间件：验证JWT令牌（接口鉴权）
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取令牌（格式：Bearer <token>）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// utils.Fail(c, errors.AUTH_ERROR, "missing authorization token")
			// 终止当前请求的后续处理流程
			c.Abort()
			return
		}

		// 解析Bearer令牌
		var tokenString string
		parts := []rune(authHeader)
		if len(parts) > 7 && string(parts[:7]) == "Bearer " {
			tokenString = string(parts[7:])
		} else {
			// utils.Fail(c, errors.AUTH_ERROR, "invalid token format (expected Bearer <token>)")
			c.Abort()
			return
		}

		// 验证令牌
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			// utils.Fail(c, errors.AUTH_ERROR, "invalid token: "+err.Error())
			c.Abort()
			return
		}

		// 将用户信息存入上下文，供后续接口使用
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		// 继续处理请求
		c.Next()
	}
}
