package install

import (
	"github.com/gofiber/fiber/v3"
)

func SetupInstallRoutes(r *fiber.App) {
	r.Get("/install", InstallPage)

	installApi := r.Group("/api/install")
	installApi.Get("/status", CheckInstallStatus)
	installApi.Post("/test-db-connection", TestDbConnectionHandler)
	installApi.Post("/install", InstallSite)
}
