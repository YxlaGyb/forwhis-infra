package main

import (
	"forestwhisperGo/admin/fix"
	"forestwhisperGo/admin/install"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"
	_ "forestwhisperGo/magic/routes/modules"
	"forestwhisperGo/magic/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(utils.Logger())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.Cors())
	router.SetTrustedProxies([]string{"127.0.0.1"})

	if !utils.ConfigExists() {
		install.SetupInstallRoutes(router)
	} else {
		router.GET("/", index)
		router.Static("/static", "./static")
		routes.SetupAllRoutes(router)
		fix.SetupFixRoutes(router)
	}

	router.Run(":2324")
}

func index(c *gin.Context) {
	middleware.JSONSuccess(c, "ok", map[string]string{"Status": "Ok"})
}

// router.Use(utils.Logging())
// gin.SetMode(gin.ReleaseMode)
