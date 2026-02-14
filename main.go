package main

import (
	"github.com/gin-gonic/gin"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// router.Use(utils.Logging())
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Use(middleware.RecoveryMiddleware())
	router.GET("/", index)
	routes.SetupAllRoutes(router)
	router.Run(":2324")
}

func index(c *gin.Context) {
	go func() {
		middleware.JSONResponse(c, 200, "ok", map[string]string{"Status": "Ok"})
	}()
}
