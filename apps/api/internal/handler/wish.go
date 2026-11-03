package handler

import (
	"strconv"
	"strings"

	"forwhis-infra/apps/api/internal/middleware"
	"forwhis-infra/apps/api/internal/service"

	"github.com/gofiber/fiber/v3"
)

func GetWishes(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	viewerID := 0
	if user := middleware.GetCurrentUser(c); user != nil {
		viewerID = user.Uid
	}
	return middleware.JSONSuccess(c, "获取成功喵", service.GetWishes(page, limit, viewerID))
}

func GetMyWishes(c fiber.Ctx) error {
	user := middleware.GetCurrentUserFromContext(c)
	return middleware.JSONSuccess(c, "获取成功喵", service.GetWishesByAuthor(user.Uid))
}

func DeleteWish(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return middleware.JSONError(c, "心愿ID无效喵", nil)
	}

	user := middleware.GetCurrentUserFromContext(c)
	if !service.DeleteWish(id, user.Uid) {
		return middleware.JSONError(c, "心愿不存在或无权删除喵", nil)
	}
	return middleware.JSONSuccess(c, "心愿已删除喵", nil)
}

func CreateWish(c fiber.Ctx) error {
	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}
	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请求数据不能为空喵", nil)
	}
	content := strings.TrimSpace(req.Content)
	if content == "" || len([]rune(content)) > 500 {
		return middleware.JSONError(c, "心愿内容需在1-500字之间喵", nil)
	}

	user := middleware.GetCurrentUser(c)
	authorID, author, mail := 0, strings.TrimSpace(req.Name), ""
	if user != nil {
		authorID, author, mail = user.Uid, user.Name, user.Mail
	}
	if author == "" {
		return middleware.JSONError(c, "请填写昵称喵", nil)
	}

	wish := service.CreateWish(content, authorID, author, mail)
	if authorID > 0 {
		service.CreateUserinfo(authorID, "createWish", "在心愿墙留言喵")
	}
	return middleware.JSONSuccess(c, "心愿已送达到星空喵", map[string]interface{}{"id": wish.ID})
}
