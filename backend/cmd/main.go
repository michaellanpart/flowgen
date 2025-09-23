package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/michaellanpart/flowgen/backend/internal/api"
	"github.com/michaellanpart/flowgen/backend/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "flowgen-backend",
			"version": "0.1.0",
		})
	})

	// Serve static files from the new frontend directory
	r.Static("/static", "../frontend")
	r.StaticFile("/", "../frontend/index.html")
	r.StaticFile("/flowchart-display.html", "../frontend/index.html")

	// API routes
	api.SetupRoutes(r)

	// Start server
	log.Printf("Starting FlowGen backend server on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
