package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"forestwhisperGo/admin/fix"
	"forestwhisperGo/admin/install"
	"forestwhisperGo/magic/middleware"
	"forestwhisperGo/magic/routes"
	_ "forestwhisperGo/magic/routes/modules"
	"forestwhisperGo/magic/utils"
	"forestwhisperGo/magic/utils/db"

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

	srv := &http.Server{
		Addr:    ":2324",
		Handler: router,
	}

	go func() {
		fmt.Println("Server starting on :2324")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	}

	db.CloseDB()
	utils.Shutdown()
	fmt.Println("Server exited")
}

func index(c *gin.Context) {
	middleware.JSONSuccess(c, "ok", map[string]string{"Status": "Ok"})
}
