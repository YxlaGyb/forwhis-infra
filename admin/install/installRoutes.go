package install

import (
	"github.com/gin-gonic/gin"
)

func SetupInstallRoutes(r *gin.Engine) {
	r.GET("/install", InstallPage)

	installApi := r.Group("/api/install")
	{
		installApi.GET("/status", CheckInstallStatus)
		installApi.POST("/test-db-connection", TestDbConnectionHandler)
		installApi.POST("/install", InstallSite)
	}
}
