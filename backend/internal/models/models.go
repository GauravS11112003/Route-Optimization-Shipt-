package models

// Order represents a delivery order
type Order struct {
	ID             string  `json:"id"`
	Lat            float64 `json:"lat"`
	Lng            float64 `json:"lng"`
	ItemCount      int     `json:"itemCount"`
	DeliveryWindow string  `json:"deliveryWindow"`
}

// Shopper represents an available delivery shopper
type Shopper struct {
	ID       string  `json:"id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Capacity int     `json:"capacity"`
}

// Assignment represents a shopper's optimized route
type Assignment struct {
	ShopperID     string   `json:"shopperId"`
	Route         []string `json:"route"`
	TotalDistance float64  `json:"totalDistance"`
}

// SampleDataResponse contains mock data for testing
type SampleDataResponse struct {
	Orders   []Order   `json:"orders"`
	Shoppers []Shopper `json:"shoppers"`
}

// OptimizeRequest contains data to be optimized
type OptimizeRequest struct {
	Orders   []Order   `json:"orders"`
	Shoppers []Shopper `json:"shoppers"`
}

// OptimizeResponse contains optimization results
type OptimizeResponse struct {
	Assignments         []Assignment `json:"assignments"`
	TotalDistanceBefore float64      `json:"totalDistanceBefore"`
	TotalDistanceAfter  float64      `json:"totalDistanceAfter"`
}

