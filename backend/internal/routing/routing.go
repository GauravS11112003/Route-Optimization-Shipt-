package routing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
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
// The geometry can be either a string (encoded) or an array of coordinates
type openRouteServiceResponse struct {
	Routes []struct {
		Summary struct {
			Distance float64 `json:"distance"` // in meters
			Duration float64 `json:"duration"` // in seconds
		} `json:"summary"`
		Geometry interface{} `json:"geometry"` // Can be string or coordinate array
	} `json:"routes"`
}

// GetRouteWithKey fetches actual driving route from OpenRouteService using provided API key
func GetRouteWithKey(fromLat, fromLng, toLat, toLng float64, apiKey string) (*RouteSegment, error) {
	// Using public OpenRouteService API
	url := "https://api.openrouteservice.org/v2/directions/driving-car/json"
	
	// Log API key status
	if apiKey != "" {
		keyPreview := apiKey
		if len(apiKey) > 10 {
			keyPreview = apiKey[:10] + "..."
		}
		println("✓ Using API Key from frontend:", keyPreview, "(", len(apiKey), "chars)")
	} else {
		println("⚠ No API Key provided from frontend")
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	// Prepare request body - request JSON format (not encoded)
	body := map[string]interface{}{
		"coordinates": [][]float64{
			{fromLng, fromLat}, // OpenRouteService uses [lng, lat] format
			{toLng, toLat},
		},
		"geometry": true, // Explicitly request geometry
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
	req.Header.Set("Authorization", apiKey) // Use provided API key
	
	// Make request with timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		// Fallback to haversine if routing service fails
		println("⚠ OpenRouteService HTTP request failed:", err.Error())
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		// Fallback if API fails
		bodyBytes, _ := io.ReadAll(resp.Body)
		println("⚠ OpenRouteService API error - Status:", resp.StatusCode)
		println("   Response body:", string(bodyBytes))
		println("   Request URL:", url)
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	println("✓ OpenRouteService API success - Status:", resp.StatusCode)
	
	// Parse response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		println("❌ Error reading response body:", err.Error())
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	// Debug: Show response sample
	if len(bodyBytes) > 200 {
		println("📦 Response sample:", string(bodyBytes[:200])+"...")
	} else {
		println("📦 Full response:", string(bodyBytes))
	}
	
	var orsResp openRouteServiceResponse
	if err := json.Unmarshal(bodyBytes, &orsResp); err != nil {
		println("❌ Error unmarshaling JSON:", err.Error())
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	if len(orsResp.Routes) == 0 {
		println("⚠️ No routes in response")
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	route := orsResp.Routes[0]
	println("📏 Route distance:", route.Summary.Distance, "meters")
	println("⏱️ Route duration:", route.Summary.Duration, "seconds")
	
	// Parse geometry - it can be a map (GeoJSON), array, or string (encoded)
	geometry := []RoutePoint{}
	
	switch geom := route.Geometry.(type) {
	case map[string]interface{}:
		// GeoJSON format: {"coordinates": [[lng, lat], ...]}
		if coords, ok := geom["coordinates"].([]interface{}); ok {
			println("🗺️ GeoJSON format detected, coordinate count:", len(coords))
			for i, coordInterface := range coords {
				if coord, ok := coordInterface.([]interface{}); ok && len(coord) >= 2 {
					lng, _ := coord[0].(float64)
					lat, _ := coord[1].(float64)
					geometry = append(geometry, RoutePoint{Lat: lat, Lng: lng})
					if i < 3 {
						println(fmt.Sprintf("   Coord %d: [%.6f, %.6f]", i, lat, lng))
					}
				}
			}
		}
	case string:
		// Encoded polyline format - decode it
		println("🗺️ Encoded polyline format:", geom[:50]+"...")
		println("🔓 Decoding polyline...")
		geometry = decodePolyline(geom)
		println("✅ Decoded to", len(geometry), "points")
		if len(geometry) > 0 {
			println(fmt.Sprintf("   First point: [%.6f, %.6f]", geometry[0].Lat, geometry[0].Lng))
			if len(geometry) > 1 {
				println(fmt.Sprintf("   Last point: [%.6f, %.6f]", geometry[len(geometry)-1].Lat, geometry[len(geometry)-1].Lng))
			}
		}
	case []interface{}:
		// Direct array of coordinates
		println("🗺️ Direct array format, coordinate count:", len(geom))
		for i, coordInterface := range geom {
			if coord, ok := coordInterface.([]interface{}); ok && len(coord) >= 2 {
				lng, _ := coord[0].(float64)
				lat, _ := coord[1].(float64)
				geometry = append(geometry, RoutePoint{Lat: lat, Lng: lng})
				if i < 3 {
					println(fmt.Sprintf("   Coord %d: [%.6f, %.6f]", i, lat, lng))
				}
			}
		}
	default:
		println("❌ Unknown geometry format:", fmt.Sprintf("%T", geom))
		return getFallbackRoute(fromLat, fromLng, toLat, toLng), nil
	}
	
	println("✅ Final geometry has", len(geometry), "points")
	
	return &RouteSegment{
		Distance: route.Summary.Distance / 1000.0, // convert meters to km
		Duration: route.Summary.Duration / 60.0,   // convert seconds to minutes
		Geometry: geometry,
	}, nil
}

// GetRoute fetches actual driving route from OpenRouteService (legacy - uses env var)
func GetRoute(fromLat, fromLng, toLat, toLng float64) (*RouteSegment, error) {
	// For backward compatibility, check environment variable
	apiKey := os.Getenv("OPENROUTE_API_KEY")
	return GetRouteWithKey(fromLat, fromLng, toLat, toLng, apiKey)
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

// decodePolyline decodes an encoded polyline string into a slice of RoutePoints
// This uses the standard Google Polyline encoding algorithm
func decodePolyline(encoded string) []RoutePoint {
	var points []RoutePoint
	index := 0
	lat := 0
	lng := 0
	
	for index < len(encoded) {
		// Decode latitude
		var shift uint = 0
		var result int = 0
		var b int
		
		for {
			if index >= len(encoded) {
				break
			}
			b = int(encoded[index]) - 63
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		
		var dlat int
		if (result & 1) != 0 {
			dlat = ^(result >> 1)
		} else {
			dlat = result >> 1
		}
		lat += dlat
		
		// Decode longitude
		shift = 0
		result = 0
		
		for {
			if index >= len(encoded) {
				break
			}
			b = int(encoded[index]) - 63
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		
		var dlng int
		if (result & 1) != 0 {
			dlng = ^(result >> 1)
		} else {
			dlng = result >> 1
		}
		lng += dlng
		
		points = append(points, RoutePoint{
			Lat: float64(lat) / 1e5,
			Lng: float64(lng) / 1e5,
		})
	}
	
	return points
}

