package optimizer

import (
	"math"
	"shipt-route-optimizer/internal/models"
	"sort"
)

// HaversineDistance calculates the distance between two lat/lng points in kilometers
func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0 // Earth radius in kilometers

	// Convert to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// OrderDistance represents an order with its distance from a point
type OrderDistance struct {
	Order    models.Order
	Distance float64
}

// Optimize assigns orders to shoppers using nearest-neighbor clustering
func Optimize(orders []models.Order, shoppers []models.Shopper) ([]models.Assignment, float64, float64) {
	if len(shoppers) == 0 || len(orders) == 0 {
		return []models.Assignment{}, 0, 0
	}

	// Track assignments per shopper
	assignments := make(map[string][]models.Order)
	for _, shopper := range shoppers {
		assignments[shopper.ID] = []models.Order{}
	}

	// Calculate total distance before optimization (random assignment)
	totalDistanceBefore := calculateRandomDistance(orders, shoppers)

	// Assign each order to nearest available shopper
	remainingOrders := make([]models.Order, len(orders))
	copy(remainingOrders, orders)

	for len(remainingOrders) > 0 {
		order := remainingOrders[0]
		remainingOrders = remainingOrders[1:]

		// Find nearest shopper with capacity
		bestShopperID := ""
		minDistance := math.MaxFloat64

		for _, shopper := range shoppers {
			if len(assignments[shopper.ID]) >= shopper.Capacity {
				continue // Shopper at capacity
			}

			distance := HaversineDistance(
				order.Lat, order.Lng,
				shopper.Lat, shopper.Lng,
			)

			if distance < minDistance {
				minDistance = distance
				bestShopperID = shopper.ID
			}
		}

		// If no shopper available, assign to first shopper (overflow)
		if bestShopperID == "" {
			bestShopperID = shoppers[0].ID
		}

		assignments[bestShopperID] = append(assignments[bestShopperID], order)
	}

	// Build optimized routes for each shopper
	result := []models.Assignment{}
	totalDistanceAfter := 0.0

	for _, shopper := range shoppers {
		shopperOrders := assignments[shopper.ID]
		if len(shopperOrders) == 0 {
			continue
		}

		// Sort orders by proximity for efficient routing (nearest neighbor)
		route := optimizeShopperRoute(shopper, shopperOrders)
		routeIDs := []string{}
		routeDistance := 0.0

		// Calculate route distance
		currentLat, currentLng := shopper.Lat, shopper.Lng
		for _, order := range route {
			routeIDs = append(routeIDs, order.ID)
			distance := HaversineDistance(currentLat, currentLng, order.Lat, order.Lng)
			routeDistance += distance
			currentLat, currentLng = order.Lat, order.Lng
		}

		result = append(result, models.Assignment{
			ShopperID:     shopper.ID,
			Route:         routeIDs,
			TotalDistance: math.Round(routeDistance*100) / 100,
		})

		totalDistanceAfter += routeDistance
	}

	return result, math.Round(totalDistanceBefore*100) / 100, math.Round(totalDistanceAfter*100) / 100
}

// optimizeShopperRoute sorts orders by nearest neighbor from shopper location
func optimizeShopperRoute(shopper models.Shopper, orders []models.Order) []models.Order {
	if len(orders) <= 1 {
		return orders
	}

	// Greedy nearest-neighbor approach
	route := []models.Order{}
	remaining := make([]models.Order, len(orders))
	copy(remaining, orders)

	currentLat, currentLng := shopper.Lat, shopper.Lng

	for len(remaining) > 0 {
		// Find nearest order
		nearestIdx := 0
		minDist := HaversineDistance(currentLat, currentLng, remaining[0].Lat, remaining[0].Lng)

		for i := 1; i < len(remaining); i++ {
			dist := HaversineDistance(currentLat, currentLng, remaining[i].Lat, remaining[i].Lng)
			if dist < minDist {
				minDist = dist
				nearestIdx = i
			}
		}

		// Add nearest order to route
		route = append(route, remaining[nearestIdx])
		currentLat, currentLng = remaining[nearestIdx].Lat, remaining[nearestIdx].Lng

		// Remove from remaining
		remaining = append(remaining[:nearestIdx], remaining[nearestIdx+1:]...)
	}

	return route
}

// calculateRandomDistance simulates random assignment for comparison
func calculateRandomDistance(orders []models.Order, shoppers []models.Shopper) float64 {
	if len(shoppers) == 0 || len(orders) == 0 {
		return 0
	}

	totalDist := 0.0
	ordersPerShopper := len(orders) / len(shoppers)
	if ordersPerShopper == 0 {
		ordersPerShopper = 1
	}

	for i, order := range orders {
		shopperIdx := i / ordersPerShopper
		if shopperIdx >= len(shoppers) {
			shopperIdx = len(shoppers) - 1
		}

		shopper := shoppers[shopperIdx]
		dist := HaversineDistance(shopper.Lat, shopper.Lng, order.Lat, order.Lng)
		totalDist += dist

		// Add distance between orders (simplified)
		if i > 0 && (i-1)/ordersPerShopper == shopperIdx {
			prevOrder := orders[i-1]
			totalDist += HaversineDistance(prevOrder.Lat, prevOrder.Lng, order.Lat, order.Lng)
		}
	}

	return totalDist
}

// SortAssignmentsByShopper sorts assignments by shopper ID for consistency
func SortAssignmentsByShopper(assignments []models.Assignment) {
	sort.Slice(assignments, func(i, j int) bool {
		return assignments[i].ShopperID < assignments[j].ShopperID
	})
}

