package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"courier-service/config"
	"courier-service/internal/database"
	"courier-service/internal/models"
	"courier-service/internal/routers"
)

func main() {
	r := gin.Default()

	// Load config
	config.LoadConfig()

	// Initialize database connection
	if err := database.InitializeDB(); err != nil {
		log.Fatal("Fatal error: Database intialization failed")
	}

	// Apply migrations
	if err := database.DB.AutoMigrate(&models.User{}, &models.DeliveryOrder{}); err != nil {
		log.Fatal("Fatal error: Migration failed")
	}

	// Load routes
	routers.GetAppRoutes(r)

	r.Run(":8080")
}
