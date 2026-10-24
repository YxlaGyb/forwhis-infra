package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forestwhisperGo/apps/admin/fix"
	"forestwhisperGo/apps/admin/install"
	"forestwhisperGo/apps/internal/middleware"
	"forestwhisperGo/apps/internal/routes"
	_ "forestwhisperGo/apps/internal/routes/modules"
	"forestwhisperGo/apps/pkg/db"
	pkg "forestwhisperGo/apps/pkg/logger"
	"forestwhisperGo/apps/pkg/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()
	app.Use(middleware.RecoveryMiddleware())
	app.Use(middleware.ProxyFix())
	app.Use(pkg.Logger())
	app.Use(middleware.Cors())

	if !utils.ConfigExists() {
		install.SetupInstallRoutes(app)
	} else {
		app.Get("/", index)
		app.Use("/static", static.New("./static"))
		routes.SetupAllRoutes(app)
		fix.SetupFixRoutes(app)
	}

	go func() {
		fmt.Println("Server starting on :2324")
		if err := app.Listen(":2324"); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	db.CloseDB()
	utils.Shutdown()
	fmt.Println("Server exited")
}

func index(c fiber.Ctx) error {
	return middleware.JSONSuccess(c, "ok", map[string]string{"status": "ok"})
}
