package models

// ShopperAnalytics contains performance metrics for a shopper
type ShopperAnalytics struct {
	ShopperID          string  `json:"shopperId"`
	OrdersAssigned     int     `json:"ordersAssigned"`
	TotalDistance      float64 `json:"totalDistance"`      // km
	TotalDuration      float64 `json:"totalDuration"`      // minutes
	CapacityUtilization float64 `json:"capacityUtilization"` // percentage
	AverageOrderDistance float64 `json:"averageOrderDistance"` // km
	EstimatedStartTime string  `json:"estimatedStartTime"`
	EstimatedEndTime   string  `json:"estimatedEndTime"`
	Efficiency         float64 `json:"efficiency"` // orders per hour
}

// OrderAnalytics contains insights about order distribution
type OrderAnalytics struct {
	TotalOrders         int     `json:"totalOrders"`
	AverageItemCount    float64 `json:"averageItemCount"`
	TotalItems          int     `json:"totalItems"`
	OrderDensity        float64 `json:"orderDensity"` // orders per sq km
	AverageDistance     float64 `json:"averageDistance"` // km
	UnassignedOrders    int     `json:"unassignedOrders"`
	TimeWindowBreakdown map[string]int `json:"timeWindowBreakdown"`
}

// SystemAnalytics contains overall system metrics
type SystemAnalytics struct {
	TotalShoppers       int     `json:"totalShoppers"`
	ActiveShoppers      int     `json:"activeShoppers"` // shoppers with assignments
	TotalOrders         int     `json:"totalOrders"`
	AssignedOrders      int     `json:"assignedOrders"`
	TotalDistance       float64 `json:"totalDistance"`
	TotalDuration       float64 `json:"totalDuration"` // minutes
	AverageEfficiency   float64 `json:"averageEfficiency"` // orders per hour
	OptimizationScore   float64 `json:"optimizationScore"` // 0-100
	EstimatedFuelCost   float64 `json:"estimatedFuelCost"` // USD
	CO2Saved            float64 `json:"co2Saved"` // kg
}

// RouteGeometry contains the actual road path
type RouteGeometry struct {
	ShopperID string      `json:"shopperId"`
	Points    [][]float64 `json:"points"` // [lat, lng] pairs
}

// AnalyticsResponse contains all analytics data
type AnalyticsResponse struct {
	System          SystemAnalytics      `json:"system"`
	Shoppers        []ShopperAnalytics   `json:"shoppers"`
	Orders          OrderAnalytics       `json:"orders"`
	RouteGeometries []RouteGeometry      `json:"routeGeometries"`
}

