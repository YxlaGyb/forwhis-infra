package main

import (
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
	router.Use(middleware.DBMiddleware())
	// router.Use(utils.Logging())
	// gin.SetMode(gin.ReleaseMode)
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.GET("/", index)
	routes.SetupAllRoutes(router)
	router.Run(":2324")
}

func index(c *gin.Context) {
	middleware.JSONSuccess(c, "ok", map[string]string{"Status": "Ok"})
}
