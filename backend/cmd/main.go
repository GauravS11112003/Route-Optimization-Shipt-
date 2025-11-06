package main

import (
	"log"
	"shipt-route-optimizer/internal/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configure CORS for frontend
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	// Register API routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/health", api.HealthCheck)
		apiGroup.GET("/sample-data", api.GetSampleData)
		apiGroup.POST("/optimize", api.OptimizeRoutes)
		apiGroup.POST("/optimize-analytics", api.OptimizeWithAnalytics)
	}

	log.Println("🚀 Shipt Route Optimizer Backend starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

