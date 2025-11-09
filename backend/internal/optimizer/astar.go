package optimizer

import (
	"container/heap"
	"math"
	"shipt-route-optimizer/internal/models"
)

// AStarNode represents a node in the A* search space
type AStarNode struct {
	orders      []models.Order // Remaining orders to visit
	route       []models.Order // Orders already visited
	currentLat  float64        // Current position latitude
	currentLng  float64        // Current position longitude
	gCost       float64        // Actual cost from start
	hCost       float64        // Heuristic cost to goal
	fCost       float64        // Total cost (g + h)
	index       int            // Index in priority queue
}

// PriorityQueue implements heap.Interface for A* nodes
type PriorityQueue []*AStarNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Lower fCost has higher priority
	return pq[i].fCost < pq[j].fCost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*AStarNode)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil  // avoid memory leak
	node.index = -1 // for safety
	*pq = old[0 : n-1]
	return node
}

// OptimizeRouteAStar uses A* algorithm to find optimal route through orders
func OptimizeRouteAStar(shopper models.Shopper, orders []models.Order) []models.Order {
	if len(orders) <= 1 {
		return orders
	}

	// For small number of orders, use exhaustive A*
	// For larger sets, use A* with beam search to limit memory
	if len(orders) <= 8 {
		return aStarExhaustive(shopper, orders)
	}
	return aStarBeamSearch(shopper, orders, 100) // Beam width of 100
}

// aStarExhaustive performs complete A* search for optimal route
func aStarExhaustive(shopper models.Shopper, orders []models.Order) []models.Order {
	// Initialize priority queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Create initial node
	initialNode := &AStarNode{
		orders:     orders,
		route:      []models.Order{},
		currentLat: shopper.Lat,
		currentLng: shopper.Lng,
		gCost:      0,
		hCost:      calculateHeuristic(shopper.Lat, shopper.Lng, orders),
		fCost:      0,
	}
	initialNode.fCost = initialNode.gCost + initialNode.hCost
	heap.Push(&pq, initialNode)

	bestRoute := orders // Fallback to original order
	bestCost := math.MaxFloat64

	// A* search
	nodesExplored := 0
	maxNodes := 50000 // Safety limit

	for pq.Len() > 0 && nodesExplored < maxNodes {
		current := heap.Pop(&pq).(*AStarNode)
		nodesExplored++

		// Goal test: all orders visited
		if len(current.orders) == 0 {
			if current.gCost < bestCost {
				bestCost = current.gCost
				bestRoute = current.route
			}
			continue
		}

		// Pruning: if current cost exceeds best, skip
		if current.gCost >= bestCost {
			continue
		}

		// Expand node: try visiting each remaining order
		for i, order := range current.orders {
			// Calculate cost to this order
			moveCost := HaversineDistance(
				current.currentLat, current.currentLng,
				order.Lat, order.Lng,
			)

			// Create new orders list without this order
			newOrders := make([]models.Order, 0, len(current.orders)-1)
			newOrders = append(newOrders, current.orders[:i]...)
			newOrders = append(newOrders, current.orders[i+1:]...)

			// Create new route with this order added
			newRoute := make([]models.Order, len(current.route), len(current.route)+1)
			copy(newRoute, current.route)
			newRoute = append(newRoute, order)

			// Create child node
			gCost := current.gCost + moveCost
			hCost := calculateHeuristic(order.Lat, order.Lng, newOrders)
			
			childNode := &AStarNode{
				orders:     newOrders,
				route:      newRoute,
				currentLat: order.Lat,
				currentLng: order.Lng,
				gCost:      gCost,
				hCost:      hCost,
				fCost:      gCost + hCost,
			}

			// Only add if promising
			if childNode.fCost < bestCost {
				heap.Push(&pq, childNode)
			}
		}
	}

	return bestRoute
}

// aStarBeamSearch uses beam search variant of A* for larger problem sizes
func aStarBeamSearch(shopper models.Shopper, orders []models.Order, beamWidth int) []models.Order {
	// Initialize with greedy nearest neighbor as baseline
	greedyRoute := optimizeShopperRoute(shopper, orders)
	greedyCost := calculateRouteCost(shopper.Lat, shopper.Lng, greedyRoute)

	// Track best solution
	bestRoute := greedyRoute
	bestCost := greedyCost

	// Beam search with A*
	currentBeam := []*AStarNode{
		{
			orders:     orders,
			route:      []models.Order{},
			currentLat: shopper.Lat,
			currentLng: shopper.Lng,
			gCost:      0,
			hCost:      calculateHeuristic(shopper.Lat, shopper.Lng, orders),
			fCost:      0,
		},
	}
	currentBeam[0].fCost = currentBeam[0].gCost + currentBeam[0].hCost

	// Iterate until all orders visited
	for len(orders) > 0 {
		nextBeam := []*AStarNode{}

		// Expand each node in current beam
		for _, node := range currentBeam {
			if len(node.orders) == 0 {
				// Complete route found
				if node.gCost < bestCost {
					bestCost = node.gCost
					bestRoute = node.route
				}
				continue
			}

			// Generate successors
			for i, order := range node.orders {
				moveCost := HaversineDistance(
					node.currentLat, node.currentLng,
					order.Lat, order.Lng,
				)

				newOrders := make([]models.Order, 0, len(node.orders)-1)
				newOrders = append(newOrders, node.orders[:i]...)
				newOrders = append(newOrders, node.orders[i+1:]...)

				newRoute := make([]models.Order, len(node.route), len(node.route)+1)
				copy(newRoute, node.route)
				newRoute = append(newRoute, order)

				gCost := node.gCost + moveCost
				hCost := calculateHeuristic(order.Lat, order.Lng, newOrders)

				childNode := &AStarNode{
					orders:     newOrders,
					route:      newRoute,
					currentLat: order.Lat,
					currentLng: order.Lng,
					gCost:      gCost,
					hCost:      hCost,
					fCost:      gCost + hCost,
				}

				nextBeam = append(nextBeam, childNode)
			}
		}

		if len(nextBeam) == 0 {
			break
		}

		// Keep only best beamWidth nodes
		if len(nextBeam) > beamWidth {
			// Sort by fCost
			for i := 0; i < len(nextBeam); i++ {
				for j := i + 1; j < len(nextBeam); j++ {
					if nextBeam[j].fCost < nextBeam[i].fCost {
						nextBeam[i], nextBeam[j] = nextBeam[j], nextBeam[i]
					}
				}
			}
			nextBeam = nextBeam[:beamWidth]
		}

		currentBeam = nextBeam
	}

	// Check remaining nodes in beam
	for _, node := range currentBeam {
		if len(node.orders) == 0 && node.gCost < bestCost {
			bestCost = node.gCost
			bestRoute = node.route
		}
	}

	return bestRoute
}

// calculateHeuristic estimates remaining cost using MST lower bound
func calculateHeuristic(currentLat, currentLng float64, orders []models.Order) float64 {
	if len(orders) == 0 {
		return 0
	}

	// Use nearest neighbor distance as simple heuristic
	// More sophisticated: use MST (Minimum Spanning Tree) of remaining orders
	minDist := math.MaxFloat64
	for _, order := range orders {
		dist := HaversineDistance(currentLat, currentLng, order.Lat, order.Lng)
		if dist < minDist {
			minDist = dist
		}
	}

	// Add MST lower bound for remaining orders
	mstCost := calculateMSTLowerBound(orders)
	
	return minDist + mstCost
}

// calculateMSTLowerBound calculates minimum spanning tree cost as lower bound
func calculateMSTLowerBound(orders []models.Order) float64 {
	if len(orders) <= 1 {
		return 0
	}

	// Prim's algorithm for MST
	visited := make(map[int]bool)
	visited[0] = true
	totalCost := 0.0

	for len(visited) < len(orders) {
		minEdge := math.MaxFloat64
		nextNode := -1

		// Find minimum edge from visited to unvisited
		for i := range orders {
			if !visited[i] {
				continue
			}
			for j := range orders {
				if visited[j] {
					continue
				}
				dist := HaversineDistance(
					orders[i].Lat, orders[i].Lng,
					orders[j].Lat, orders[j].Lng,
				)
				if dist < minEdge {
					minEdge = dist
					nextNode = j
				}
			}
		}

		if nextNode == -1 {
			break
		}

		visited[nextNode] = true
		totalCost += minEdge
	}

	return totalCost
}

// calculateRouteCost calculates total distance of a route
func calculateRouteCost(startLat, startLng float64, route []models.Order) float64 {
	if len(route) == 0 {
		return 0
	}

	totalCost := 0.0
	currentLat, currentLng := startLat, startLng

	for _, order := range route {
		dist := HaversineDistance(currentLat, currentLng, order.Lat, order.Lng)
		totalCost += dist
		currentLat, currentLng = order.Lat, order.Lng
	}

	return totalCost
}

// OptimizeAStar performs full optimization using A* for route planning
func OptimizeAStar(orders []models.Order, shoppers []models.Shopper) ([]models.Assignment, float64, float64) {
	if len(shoppers) == 0 || len(orders) == 0 {
		return []models.Assignment{}, 0, 0
	}

	// Track assignments per shopper
	assignments := make(map[string][]models.Order)
	for _, shopper := range shoppers {
		assignments[shopper.ID] = []models.Order{}
	}

	// Calculate baseline distance
	totalDistanceBefore := calculateRandomDistance(orders, shoppers)

	// Assign each order to nearest available shopper (greedy assignment)
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
				continue
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

		if bestShopperID == "" {
			bestShopperID = shoppers[0].ID
		}

		assignments[bestShopperID] = append(assignments[bestShopperID], order)
	}

	// Build optimized routes using A* for each shopper
	result := []models.Assignment{}
	totalDistanceAfter := 0.0

	for _, shopper := range shoppers {
		shopperOrders := assignments[shopper.ID]
		if len(shopperOrders) == 0 {
			continue
		}

		// Use A* to optimize route sequence
		route := OptimizeRouteAStar(shopper, shopperOrders)
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

