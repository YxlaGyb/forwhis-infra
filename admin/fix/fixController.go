package fix

import (
	"forestwhisperGo/magic/middleware"

	"github.com/gin-gonic/gin"
)

func CheckTablesHandler(c *gin.Context) {
	if !CheckFixToolAllowed() {
		middleware.JSONError(c, "修复工具未启用，请检查数据库连接", nil)
		return
	}

	tables := CheckTables()
	middleware.JSONSuccess(c, "获取成功", gin.H{
		"tables": tables,
	})
}

func FixTablesHandler(c *gin.Context) {
	if !CheckFixToolAllowed() {
		middleware.JSONError(c, "修复工具未启用，请检查数据库连接", nil)
		return
	}

	fixedCount, results := FixTables()
	middleware.JSONSuccess(c, "修复完成", gin.H{
		"fixedCount": fixedCount,
		"results":    results,
	})
}

func FixSingleTableHandler(c *gin.Context) {
	if !CheckFixToolAllowed() {
		middleware.JSONError(c, "修复工具未启用，请检查数据库连接", nil)
		return
	}

	var req struct {
		TableName string `json:"tableName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.JSONError(c, "请提供表名", nil)
		return
	}

	if err := FixSingleTable(req.TableName); err != nil {
		middleware.JSONError(c, "修复失败: "+err.Error(), nil)
		return
	}

	middleware.JSONSuccess(c, "修复成功", gin.H{
		"tableName": req.TableName,
	})
}

func GetTableInfoHandler(c *gin.Context) {
	if !CheckFixToolAllowed() {
		middleware.JSONError(c, "修复工具未启用，请检查数据库连接", nil)
		return
	}

	tableName := c.Query("table")
	if tableName == "" {
		middleware.JSONError(c, "请提供表名", nil)
		return
	}

	info, err := GetTableInfo(tableName)
	if err != nil {
		middleware.JSONError(c, "获取表信息失败: "+err.Error(), nil)
		return
	}

	middleware.JSONSuccess(c, "获取成功", info)
}

func CreateIndexesHandler(c *gin.Context) {
	if !CheckFixToolAllowed() {
		middleware.JSONError(c, "修复工具未启用，请检查数据库连接", nil)
		return
	}

	createdCount, results := CreateMissingIndexes()
	middleware.JSONSuccess(c, "索引创建完成", gin.H{
		"createdCount": createdCount,
		"results":      results,
	})
}

func FixStatusHandler(c *gin.Context) {
	middleware.JSONSuccess(c, "获取成功", gin.H{
		"allowed": CheckFixToolAllowed(),
	})
}
