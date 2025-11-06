package data

import (
	"fmt"
	"math/rand"
	"shipt-route-optimizer/internal/models"
	"time"
)

// GenerateSampleData creates mock orders and shoppers around Birmingham, AL (Shipt HQ)
func GenerateSampleData() models.SampleDataResponse {
	rand.Seed(time.Now().UnixNano())

	// Birmingham, AL coordinates (Shipt headquarters area)
	centerLat := 33.5186
	centerLng := -86.8104
	radius := 0.15 // ~10 mile radius

	// Generate 5 shoppers
	shoppers := []models.Shopper{}
	for i := 1; i <= 5; i++ {
		shoppers = append(shoppers, models.Shopper{
			ID:       fmt.Sprintf("S%d", i),
			Lat:      centerLat + (rand.Float64()-0.5)*radius,
			Lng:      centerLng + (rand.Float64()-0.5)*radius,
			Capacity: rand.Intn(3) + 3, // 3-5 orders per shopper
		})
	}

	// Generate 20 orders
	orders := []models.Order{}
	deliveryWindows := []string{"9-11 AM", "11 AM-1 PM", "1-3 PM", "3-5 PM", "5-7 PM"}

	for i := 1; i <= 20; i++ {
		orders = append(orders, models.Order{
			ID:             fmt.Sprintf("O%d", i),
			Lat:            centerLat + (rand.Float64()-0.5)*radius,
			Lng:            centerLng + (rand.Float64()-0.5)*radius,
			ItemCount:      rand.Intn(30) + 5, // 5-35 items
			DeliveryWindow: deliveryWindows[rand.Intn(len(deliveryWindows))],
		})
	}

	return models.SampleDataResponse{
		Orders:   orders,
		Shoppers: shoppers,
	}
}

