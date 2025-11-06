package api

import (
	"net/http"
	"shipt-route-optimizer/internal/data"
	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/optimizer"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns API health status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "shipt-route-optimizer",
	})
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

	// Run optimization with analytics
	optimizeResponse, analyticsResponse := optimizer.OptimizeWithAnalytics(
		req.Orders,
		req.Shoppers,
		req.UseRealRoutes,
	)

	// Combine both responses
	response := gin.H{
		"optimization": optimizeResponse,
		"analytics":    analyticsResponse,
	}

	c.JSON(http.StatusOK, response)
}

