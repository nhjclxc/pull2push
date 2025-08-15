package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// GlobalPanicRecovery 捕获所有 panic
func GlobalPanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errMsg := fmt.Sprintf("Panic: %v\n%s", r, debug.Stack())
				fmt.Println(errMsg) // 可以换成日志输出 logger.Error()

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "服务器内部错误，请稍后再试",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
