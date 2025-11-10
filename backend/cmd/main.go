package main

import (
	"log"
	"path/filepath"
	"shipt-route-optimizer/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file - try multiple locations
	envLocations := []string{
		".env",                    // Current directory
		"../.env",                 // Parent directory (if running from cmd/)
		filepath.Join("backend", ".env"), // From project root
	}
	
	for _, envPath := range envLocations {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("Loaded .env file from: %s\n", envPath)
			break
		}
	}
	
	r := gin.Default()

	// Configure CORS for frontend
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	// Register API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", api.HealthCheck)
		apiGroup.GET("/test-routing", api.TestRouting)
		apiGroup.GET("/sample-data", api.GetSampleData)
		apiGroup.POST("/optimize", api.OptimizeRoutes)
		apiGroup.POST("/optimize-analytics", api.OptimizeWithAnalytics)
	}

	log.Println("Shipt Route Optimizer Backend starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

