package api

import (
	"net/http"
	"os"
	"shipt-route-optimizer/internal/data"
	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/optimizer"
	"shipt-route-optimizer/internal/routing"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns API health status
func HealthCheck(c *gin.Context) {
	apiKeySet := os.Getenv("OPENROUTE_API_KEY") != ""
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "shipt-route-optimizer",
		"apiKeySet": apiKeySet,
	})
}

// TestRouting tests the OpenRouteService API
func TestRouting(c *gin.Context) {
	// Test route from Birmingham coordinates  
	segment, err := routing.GetRoute(33.5200, -86.8100, 33.5186, -86.8104)
	
	result := gin.H{
		"error":          nil,
		"pointCount":     0,
		"distance":       0.0,
		"duration":       0.0,
		"apiKeySet":      os.Getenv("OPENROUTE_API_KEY") != "",
		"apiKeyLength":   len(os.Getenv("OPENROUTE_API_KEY")),
		"usingFallback":  false,
	}
	
	if err != nil {
		result["error"] = err.Error()
	}
	
	if segment != nil {
		result["pointCount"] = len(segment.Geometry)
		result["distance"] = segment.Distance
		result["duration"] = segment.Duration
		// If only 2 points, it's using fallback straight line
		result["usingFallback"] = len(segment.Geometry) == 2
		if len(segment.Geometry) <= 5 {
			result["geometry"] = segment.Geometry
		} else {
			result["geometrySample"] = gin.H{
				"first": segment.Geometry[0],
				"last":  segment.Geometry[len(segment.Geometry)-1],
				"total": len(segment.Geometry),
			}
		}
	}
	
	c.JSON(http.StatusOK, result)
}

// GetSampleData returns mock orders and shoppers
func GetSampleData(c *gin.Context) {
	sampleData := data.GenerateSampleData()
	c.JSON(http.StatusOK, sampleData)
}

// OptimizeRoutes assigns orders to shoppers and optimizes routes
func OptimizeRoutes(c *gin.Context) {
	var req models.OptimizeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.Orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No orders provided"})
		return
	}

	if len(req.Shoppers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No shoppers provided"})
		return
	}

	// Run optimization algorithm
	assignments, totalBefore, totalAfter := optimizer.Optimize(req.Orders, req.Shoppers)
	optimizer.SortAssignmentsByShopper(assignments)

	response := models.OptimizeResponse{
		Assignments:         assignments,
		TotalDistanceBefore: totalBefore,
		TotalDistanceAfter:  totalAfter,
	}

	c.JSON(http.StatusOK, response)
}

// OptimizeWithAnalytics performs optimization and returns detailed analytics
func OptimizeWithAnalytics(c *gin.Context) {
	var req struct {
		Orders        []models.Order   `json:"orders"`
		Shoppers      []models.Shopper `json:"shoppers"`
		UseRealRoutes bool             `json:"useRealRoutes"`
		Algorithm     string           `json:"algorithm"` // "nearest-neighbor" or "astar"
		ApiKey        string           `json:"apiKey"`    // OpenRouteService API key from frontend
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if len(req.Orders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No orders provided"})
		return
	}

	if len(req.Shoppers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No shoppers provided"})
		return
	}

	// Default to nearest-neighbor if not specified
	if req.Algorithm == "" {
		req.Algorithm = "nearest-neighbor"
	}


	// Run optimization with analytics
	optimizeResponse, analyticsResponse := optimizer.OptimizeWithAnalytics(
		req.Orders,
		req.Shoppers,
		req.UseRealRoutes,
		req.Algorithm,
		req.ApiKey, // Pass API key to optimizer
	)

	// Combine both responses
	response := gin.H{
		"optimization": optimizeResponse,
		"analytics":    analyticsResponse,
		"algorithm":    req.Algorithm,
	}

	c.JSON(http.StatusOK, response)
}

