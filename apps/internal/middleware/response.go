package middleware

import (
	"fmt"
	"runtime/debug"

	"forestwhisperGo/apps/internal/response"

	"github.com/gofiber/fiber/v3"
)

func Success(msg string, data interface{}) response.Response {
	return response.Response{Code: 200, Msg: msg, Data: data}
}

func Error(msg string, data interface{}) response.Response {
	return response.Response{Code: 233, Msg: msg, Data: data}
}

func SystemError(msg string, data interface{}) response.Response {
	return response.Response{Code: 500, Msg: msg, Data: data}
}

func JSONSuccess(c fiber.Ctx, msg string, data interface{}) error {
	return c.Status(200).JSON(Success(msg, data))
}

func JSONError(c fiber.Ctx, msg string, data interface{}) error {
	return c.Status(233).JSON(Error(msg, data))
}

func JSONSystemError(c fiber.Ctx, msg string, data interface{}) error {
	return c.Status(500).JSON(SystemError(msg, data))
}

func RecoveryMiddleware() fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				fmt.Printf("[PANIC] %v\n%s\n", recovered, debug.Stack())
				err = JSONSystemError(c, "服务器繁忙, 请稍后再试喵", nil)
			}
		}()
		return c.Next()
	}
}
