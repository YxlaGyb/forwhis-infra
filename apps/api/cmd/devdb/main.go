package main

import (
	"fmt"
	"log"

	"forwhis-infra/apps/api/internal/models"
	"forwhis-infra/apps/api/pkg/db"
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
}
