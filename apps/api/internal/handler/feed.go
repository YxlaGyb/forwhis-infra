package handler

import (
	"strconv"

	"forwhis-infra/apps/api/internal/feed"
	"forwhis-infra/apps/api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func GetFeed(c fiber.Ctx) error {
	mode := c.Query("mode", feed.ModeRecommend)
	if mode != feed.ModeRecommend && mode != feed.ModeLatest {
		return middleware.JSONError(c, "mode 只能是 recommend 或 latest 喵", nil)
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	result, err := feed.GetFeed(mode, c.Query("cursor"), limit)
	if err != nil {
		if err == feed.ErrInvalidCursor {
			return middleware.JSONError(c, "游标无效喵", nil)
		}
		return middleware.JSONSystemError(c, "获取推荐失败喵", nil)
	}
	return middleware.JSONSuccess(c, "获取成功喵", result)
}
