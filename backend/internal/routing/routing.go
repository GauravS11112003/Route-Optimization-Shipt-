package routing

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"time"
)

// RoutePoint represents a coordinate
type RoutePoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// RouteSegment contains route information between two points
type RouteSegment struct {
	Distance float64       `json:"distance"` // in kilometers
	Duration float64       `json:"duration"` // in minutes
	Geometry []RoutePoint  `json:"geometry"` // actual road path
}

// OpenRouteServiceResponse represents the API response
type openRouteServiceResponse struct {
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"` // in meters
			Duration float64 `json:"duration"` // in seconds
		} `json:"summary"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

// GetRoute fetches actual driving route from OpenRouteService
func GetRoute(fromLat, fromLng, toLat, toLng float64) (*RouteSegment, error) {
	// Using public OpenRouteService API (free tier, no key required for basic usage)
	url := "https://api.openrouteservice.org/v2/directions/driving-car"
	
	// Prepare request body
	body := map[string]interface{}{
		"coordinates": [][]float64{
			{fromLng, fromLat}, // OpenRouteService uses [lng, lat] format
			{toLng, toLat},
		},
	}
	
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	
	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	// Note: For production, use an API key: req.Header.Set("Authorization", "YOUR_API_KEY")
	
	// Make request with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Fallback to haversine if routing service fails
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// Fallback if API fails
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	// Parse response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	var orsResp openRouteServiceResponse
	if err := json.Unmarshal(bodyBytes, &orsResp); err != nil {
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	if len(orsResp.Routes) == 0 {
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	route := orsResp.Routes[0]
	
	// Convert geometry to RoutePoints
	geometry := []RoutePoint{}
	for _, coord := range route.Geometry.Coordinates {
		if len(coord) >= 2 {
			geometry = append(geometry, RoutePoint{
				Lat: coord[1],
				Lng: coord[0],
			})
		}
	}
	
	return &RouteSegment{
		Distance: route.Summary.Distance / 1000.0, // convert meters to km
		Duration: route.Summary.Duration / 60.0,   // convert seconds to minutes
		Geometry: geometry,
	}, nil
}

// getFallbackRoute creates a simple straight line route using haversine
func getFallbackRoute(fromLat, fromLng, toLat, toLng float64) *RouteSegment {
	distance := haversineDistance(fromLat, fromLng, toLat, toLng)
	
	// Estimate duration assuming average speed of 40 km/h in city
	duration := (distance / 40.0) * 60.0 // in minutes
	
	return &RouteSegment{
		Distance: distance,
		Duration: duration,
		Geometry: []RoutePoint{
			{Lat: fromLat, Lng: fromLng},
			{Lat: toLat, Lng: toLng},
		},
	}
}

// haversineDistance calculates distance between two points in kilometers
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371.0
	
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLng := (lng2 - lng1) * math.Pi / 180
	
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

// GetMultiPointRoute calculates route through multiple points
func GetMultiPointRoute(points []RoutePoint) (*RouteSegment, error) {
	if len(points) < 2 {
		return &RouteSegment{}, nil
	}
	
	totalDistance := 0.0
	totalDuration := 0.0
	allGeometry := []RoutePoint{}
	
	// Calculate route between consecutive points
	for i := 0; i < len(points)-1; i++ {
		segment, err := GetRoute(
			points[i].Lat, points[i].Lng,
			points[i+1].Lat, points[i+1].Lng,
		)
		if err != nil {
			continue
		}
		
		totalDistance += segment.Distance
		totalDuration += segment.Duration
		allGeometry = append(allGeometry, segment.Geometry...)
	}
	
	return &RouteSegment{
		Distance: totalDistance,
		Duration: totalDuration,
		Geometry: allGeometry,
	}, nil
}

// BatchGetRoutes fetches multiple routes in parallel (with rate limiting)
func BatchGetRoutes(pairs [][]RoutePoint) ([]*RouteSegment, error) {
	results := make([]*RouteSegment, len(pairs))
	
	// For simplicity, process sequentially with small delay
	for i, pair := range pairs {
		if len(pair) != 2 {
			results[i] = getFallbackRoute(0, 0, 0, 0)
			continue
		}
		
		segment, err := GetRoute(pair[0].Lat, pair[0].Lng, pair[1].Lat, pair[1].Lng)
		if err != nil {
			results[i] = getFallbackRoute(pair[0].Lat, pair[0].Lng, pair[1].Lat, pair[1].Lng)
		} else {
			results[i] = segment
		}
		
		// Small delay to respect rate limits (5 requests per second)
		time.Sleep(200 * time.Millisecond)
	}
	
	return results, nil
}

