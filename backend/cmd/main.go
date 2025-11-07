package main

import (
	"log"
	"os"
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
	
	envLoaded := false
	for _, envPath := range envLocations {
		if err := godotenv.Load(envPath); err == nil {
			log.Printf("✓ Loaded .env file from: %s\n", envPath)
			envLoaded = true
			break
		}
	}
	
	if !envLoaded {
		log.Println("⚠ Warning: .env file not found, using system environment variables")
	}
	
	// Log API key status at startup
	if apiKey := os.Getenv("OPENROUTE_API_KEY"); apiKey != "" {
		log.Printf("✓ OpenRouteService API Key loaded (%d chars)\n", len(apiKey))
	} else {
		log.Println("⚠ WARNING: OPENROUTE_API_KEY not found in environment!")
		log.Println("   To enable real routing, add your OpenRouteService API key to .env file")
	}
	
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
		apiGroup.GET("/test-routing", api.TestRouting)
		apiGroup.GET("/sample-data", api.GetSampleData)
		apiGroup.POST("/optimize", api.OptimizeRoutes)
		apiGroup.POST("/optimize-analytics", api.OptimizeWithAnalytics)
	}

	log.Println("🚀 Shipt Route Optimizer Backend starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

