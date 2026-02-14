package middleware

import (
	"github.com/gin-gonic/gin"
	"forestwhisperGo/magic/models/types"
)


func Success(msg string, data interface{}) types.Response {
	return types.Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	}
}

func Error(msg string, data interface{}) types.Response {
	return types.Response{
		Code: 233,
		Msg:  msg,
		Data: data,
	}
}

func SystemError(msg string, data interface{}) types.Response {
	return types.Response{
		Code: 500,
		Msg:  msg,
		Data: data,
	}
}

func JSONResponse(c *gin.Context, code int, msg string, data interface{}) {
    resp := types.Response{
        Code: code,
        Msg:  msg,
        Data: data,
    }
    c.JSON(code, resp)
}

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		defer func() {
			if err := recover(); err != nil {
				JSONResponse(c, 500, "服务器繁忙, 请稍后再试喵", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
