package middleware

import (
	"bytes"
	"io/ioutil"
	"runtime"
	"time"

	"chain/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// GinMiddleware 是 Gin 的日志中间件
func GinLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 获取请求 body（需要读出来再放回去）
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		//  避免打印大请求体
		bodySnippet := string(bodyBytes)
		if len(bodyBytes) > 1024 {
			bodySnippet = bodySnippet[:1024] + "...(truncated)"
		}
		// 把 body 放回 request，后续 handler 才能继续读取
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// 处理请求
		c.Next()

		// 请求处理结束后记录日志
		end := time.Now()
		latency := end.Sub(start)

		// 获取状态码和请求信息
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// 记录请求日志
		log.Printf("HTTP request | status_code: %d, latency: %v, client_ip: %s, method: %s, path: %s, user_agent: %s, query_params: %s, body_params: %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
			userAgent,
			c.Request.URL.RawQuery,
			string(bodyBytes),
		)
	}
}

// GinRecoveryWithLogger 是带有日志记录的恢复中间件
func GinRecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 错误
				// 创建一个缓冲区来存储堆栈信息
				buf := make([]byte, 4096) // 4KB 的缓冲区通常足够
				// 调用 Stack 函数，第二个参数 true 表示打印所有协程的堆栈，false 表示只打印当前协程
				n := runtime.Stack(buf, false)

				// 将堆栈信息转换为字符串 (注意：只取有效长度 n)
				stack := string(buf[:n])
				log.Panicf("Panic recovered | error: %v, method: %s, path: %s, client_ip: %s, stack: %s",
					err,
					c.Request.Method,
					c.Request.URL.Path,
					c.ClientIP(),
					stack) // 这样你就能在日志里看到具体的报错行数了

				utils.Fail(c, 500, "system error")
				c.Abort()
			}
		}()
		c.Next()
	}
}
