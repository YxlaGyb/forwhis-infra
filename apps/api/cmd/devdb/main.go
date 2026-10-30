package main

import (
	"fmt"
	"log"
	"os"

	"forestwhisperGo/apps/api/internal/models"
	"forestwhisperGo/apps/api/internal/service"
	rbacSvc "forestwhisperGo/apps/api/internal/service/rbac"
	"forestwhisperGo/apps/api/pkg/constants"
	"forestwhisperGo/apps/api/pkg/db"
	"forestwhisperGo/apps/api/pkg/utils"
)

func main() {
	database := db.GetDB()
	if database == nil {
		log.Fatal("Database connection failed: Please check the [db] configuration in the `config.toml` file located in the root directory.")
	}

	all := models.All()
	if err := database.AutoMigrate(all...); err != nil {
		log.Fatalf("Database table synchronization failed: %v", err)
	}
	fmt.Printf("Database tables synchronized (%d tables)\n", len(all))

	if err := ensureDevAdmin(); err != nil {
		log.Fatalf("Development admin initialization failed: %v", err)
	}
}

func ensureDevAdmin() error {
	name := env("DEV_ADMIN_NAME", "admin")
	email := env("DEV_ADMIN_EMAIL", "admin@example.com")
	password := env("DEV_ADMIN_PASSWORD", "admin123456")

	if !utils.IsValidName(name) || !utils.IsValidEmail(email) || !utils.IsValidPassword(password) {
		return fmt.Errorf("DEV_ADMIN_NAME/EMAIL/PASSWORD format is invalid (password must be 8-16 characters and include both letters and numbers)")
	}

	if user := service.GetUserByEmail(email); user != nil {
		if !rbacSvc.AssignRoleToUser(user.Uid, string(constants.RoleSuperadmin), 0) {
			return fmt.Errorf("Failed to assign superadmin role to existing user %s", email)
		}
		fmt.Printf("Development admin already exists: %s\n", email)
		return nil
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %v", err)
	}
	if _, err := service.CreateUser(name, email, hash, "", string(constants.RoleSuperadmin), false); err != nil {
		return fmt.Errorf("Failed to create user: %v", err)
	}
	fmt.Printf("Development admin created: %s / %s\n", email, password)
	return nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
