package main

import (
	"log"
	"net/http"

	"github.com/SavanRajyaguru/ecommerce-go-notification-service/config"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "notification-service"})
	})

	// Admin endpoints could be added here to view logs from Mongo

	port := config.AppConfig.AppPort
	if port == "" {
		port = "8085"
	}

	log.Printf("Starting Admin API on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start API: %v", err)
	}
}
