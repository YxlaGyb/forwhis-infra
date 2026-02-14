package routes

import (
	"github.com/gin-gonic/gin"
	"forestwhisperGo/magic/routes/modules"
)

var routeRegistry = make([]func(*gin.Engine), 0)

func Register(fn func(*gin.Engine)) {
	routeRegistry = append(routeRegistry, fn)
}

func SetupAllRoutes(r *gin.Engine) {
	for _, fn := range routeRegistry {
		fn(r)
	}
}
