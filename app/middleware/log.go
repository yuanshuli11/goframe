package middleware

import (
	"github.com/gin-gonic/gin"
)

//RequestInfo 生成access日志中间件
func RequestInfo() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
