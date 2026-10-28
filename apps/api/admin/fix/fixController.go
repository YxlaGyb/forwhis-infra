package fix

import (
	"forestwhisperGo/apps/api/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func CheckTablesHandler(c fiber.Ctx) error {
	if !CheckFixToolAllowed() {
		return middleware.JSONError(c, "修复工具未启用, 请检查数据库连接", nil)
	}

	tables := CheckTables()
	return middleware.JSONSuccess(c, "获取成功", fiber.Map{
		"tables": tables,
	})

}

func FixTablesHandler(c fiber.Ctx) error {
	if !CheckFixToolAllowed() {
		return middleware.JSONError(c, "修复工具未启用, 请检查数据库连接", nil)
	}

	fixedCount, results := FixTables()
	return middleware.JSONSuccess(c, "修复完成", fiber.Map{
		"fixedCount": fixedCount,
		"results":    results,
	})

}

func FixSingleTableHandler(c fiber.Ctx) error {
	if !CheckFixToolAllowed() {
		return middleware.JSONError(c, "修复工具未启用, 请检查数据库连接", nil)
	}

	var req struct {
		TableName string `json:"tableName" binding:"required"`
	}

	if err := c.Bind().JSON(&req); err != nil {
		return middleware.JSONError(c, "请提供表名", nil)
	}

	if err := FixSingleTable(req.TableName); err != nil {
		return middleware.JSONError(c, "修复失败: "+err.Error(), nil)
	}
	return middleware.JSONSuccess(c, "修复成功", fiber.Map{
		"tableName": req.TableName,
	})

}

func GetTableInfoHandler(c fiber.Ctx) error {
	if !CheckFixToolAllowed() {
		return middleware.JSONError(c, "修复工具未启用, 请检查数据库连接", nil)
	}

	tableName := c.Query("table")
	if tableName == "" {
		return middleware.JSONError(c, "请提供表名", nil)
	}

	info, err := GetTableInfo(tableName)
	if err != nil {
		return middleware.JSONError(c, "获取表信息失败: "+err.Error(), nil)
	}
	return middleware.JSONSuccess(c, "获取成功", info)
}

func FixStatusHandler(c fiber.Ctx) error {
	return middleware.JSONSuccess(c, "获取成功", fiber.Map{
		"allowed": CheckFixToolAllowed(),
	})

}
