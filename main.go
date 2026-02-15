package main

import (
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
	router.SetTrustedProxies([]string{"127.0.0.1"})

	if !utils.ConfigExists() {
		install.SetupInstallRoutes(router)
	} else {
		router.GET("/", index)
		routes.SetupAllRoutes(router)
	}

	router.Run(":2324")
}

func index(c *gin.Context) {
	middleware.JSONSuccess(c, "ok", map[string]string{"Status": "Ok"})
}

// router.Use(utils.Logging())
// gin.SetMode(gin.ReleaseMode)
