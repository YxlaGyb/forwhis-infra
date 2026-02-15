package middleware

import (
	"fmt"
	"runtime/debug"

	"forestwhisperGo/magic/models/types"

	"github.com/gin-gonic/gin"
)

func Success(msg string, data interface{}) types.Response {
	return types.Response{Code: 200, Msg: msg, Data: data}
}

func Error(msg string, data interface{}) types.Response {
	return types.Response{Code: 233, Msg: msg, Data: data}
}

func SystemError(msg string, data interface{}) types.Response {
	return types.Response{Code: 500, Msg: msg, Data: data}
}

func JSONSuccess(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, Success(msg, data))
}

func JSONError(c *gin.Context, msg string, data interface{}) {
	c.JSON(233, Error(msg, data))
}

func JSONSystemError(c *gin.Context, msg string, data interface{}) {
	c.JSON(500, SystemError(msg, data))
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] %v\n%s\n", err, debug.Stack())
				JSONSystemError(c, "服务器繁忙, 请稍后再试喵", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
