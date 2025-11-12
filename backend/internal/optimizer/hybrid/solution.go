package hybrid

import (
	"math"
	"math/rand"
	"sort"

	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/optimizer"
)

type solution struct {
	routes         [][]int
	routeDistances []float64
	totalDistance  float64
	temperature    float64
}

func newSolution(shoppersCount int) *solution {
	routes := make([][]int, shoppersCount)
	routeDistances := make([]float64, shoppersCount)
	return &solution{
		routes:         routes,
		routeDistances: routeDistances,
		totalDistance:  0,
		temperature:    1.0,
	}
}

func (s *solution) clone() *solution {
	copyRoutes := make([][]int, len(s.routes))
	for i, route := range s.routes {
		if len(route) == 0 {
			continue
		}
		cp := make([]int, len(route))
		copy(cp, route)
		copyRoutes[i] = cp
	}
	copyDistances := make([]float64, len(s.routeDistances))
	copy(copyDistances, s.routeDistances)
	return &solution{
		routes:         copyRoutes,
		routeDistances: copyDistances,
		totalDistance:  s.totalDistance,
		temperature:    s.temperature,
	}
}

func (s *solution) recomputeTotals(cache *distanceCache) {
	total := 0.0
	for shopperIdx := range s.routes {
		s.routeDistances[shopperIdx] = cache.routeDistance(shopperIdx, s.routes[shopperIdx])
		total += s.routeDistances[shopperIdx]
	}
	s.totalDistance = total
}

func (s *solution) orderCount() int {
	total := 0
	for _, route := range s.routes {
		total += len(route)
	}
	return total
}

func (s *solution) toAssignments(
	orders []models.Order,
	shoppers []models.Shopper,
	cache *distanceCache,
) []models.Assignment {
	assignments := []models.Assignment{}
	for shopperIdx, route := range s.routes {
		if len(route) == 0 {
			continue
		}
		orderIDs := make([]string, len(route))
		for i, orderIdx := range route {
			orderIDs[i] = orders[orderIdx].ID
		}
		assignments = append(assignments, models.Assignment{
			ShopperID:     shoppers[shopperIdx].ID,
			Route:         orderIDs,
			TotalDistance: math.Round(cache.routeDistance(shopperIdx, route)*100) / 100,
		})
	}
	return assignments
}

func buildInitialSolution(cache *distanceCache, opts normalizedOptions, rng *rand.Rand) *solution {
	result := newSolution(len(cache.shoppers))
	orderIndices := rng.Perm(len(cache.orders))
	rclSize := opts.rclSize
	if rclSize < 1 {
		rclSize = 1
	}
	loads := make([]int, len(cache.shoppers))

	for _, orderIdx := range orderIndices {
		type candidate struct {
			shopper int
			dist    float64
		}

		candidates := make([]candidate, len(cache.shoppers))
		for shopperIdx := range cache.shoppers {
			candidates[shopperIdx] = candidate{
				shopper: shopperIdx,
				dist:    cache.shopperToOrder[shopperIdx][orderIdx],
			}
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].dist < candidates[j].dist
		})

		rclLimit := rclSize
		if rclLimit > len(candidates) {
			rclLimit = len(candidates)
		}
		selectedShopper := -1
		for i := 0; i < rclLimit; i++ {
			c := candidates[i]
			if cache.hasCapacity(c.shopper, loads[c.shopper]) {
				selectedShopper = c.shopper
				break
			}
		}

		if selectedShopper == -1 {
			for _, c := range candidates {
				if cache.hasCapacity(c.shopper, loads[c.shopper]) {
					selectedShopper = c.shopper
					break
				}
			}
		}

		if selectedShopper == -1 {
			selectedShopper = candidates[0].shopper
		}

		result.routes[selectedShopper] = append(result.routes[selectedShopper], orderIdx)
		loads[selectedShopper]++
	}

	// Shuffle each route using randomized NN and compute distances.
	for shopperIdx := range result.routes {
		if len(result.routes[shopperIdx]) <= 1 {
			result.routeDistances[shopperIdx] = cache.routeDistance(shopperIdx, result.routes[shopperIdx])
			continue
		}
		result.routes[shopperIdx] = randomizedNearestNeighbor(cache, shopperIdx, result.routes[shopperIdx], rng)
		result.routeDistances[shopperIdx] = cache.routeDistance(shopperIdx, result.routes[shopperIdx])
	}

	total := 0.0
	for _, dist := range result.routeDistances {
		total += dist
	}
	result.totalDistance = total
	result.temperature = math.Max(total*0.05, 1.0)

	return result
}

func randomizedNearestNeighbor(cache *distanceCache, shopperIdx int, orders []int, rng *rand.Rand) []int {
	remaining := make([]int, len(orders))
	copy(remaining, orders)

	route := make([]int, 0, len(orders))
	currentIndex := -1

	for len(remaining) > 0 {
		type neighbor struct {
			orderIdx int
			dist     float64
		}
		candidates := make([]neighbor, len(remaining))
		for i, orderIdx := range remaining {
			var dist float64
			if currentIndex == -1 {
				dist = cache.shopperToOrder[shopperIdx][orderIdx]
			} else {
				dist = cache.orderToOrder[currentIndex][orderIdx]
			}
			// Inject slight randomness to promote exploration.
			dist *= 1 + rng.Float64()*0.1
			candidates[i] = neighbor{orderIdx: orderIdx, dist: dist}
		}
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].dist < candidates[j].dist })
		pick := candidates[0].orderIdx
		route = append(route, pick)
		currentIndex = pick

		for idx := range remaining {
			if remaining[idx] == pick {
				remaining = append(remaining[:idx], remaining[idx+1:]...)
				break
			}
		}
	}
	return route
}

func runLocalSearch(base *solution, cache *distanceCache, opts normalizedOptions, rng *rand.Rand) (*solution, int) {
	best := base.clone()
	temperature := best.temperature
	if temperature <= 0 {
		temperature = math.Max(best.totalDistance*0.05, 1.0)
	}

	improvements := 0

	for iter := 0; iter < opts.localSearch; iter++ {
		neighbor := best.clone()
		orderCount := neighbor.orderCount()
		removeCount := int(math.Ceil(opts.destroyRate * float64(orderCount)))
		if removeCount < 1 {
			removeCount = 1
		}
		if removeCount > orderCount {
			removeCount = orderCount
		}
		removed := neighbor.destroy(removeCount, rng)
		neighbor.repair(removed, cache, opts, rng)
		neighbor.recomputeTotals(cache)

		delta := neighbor.totalDistance - best.totalDistance
		accept := false
		if delta < 0 {
			best = neighbor
			improvements++
			accept = true
		} else {
			threshold := math.Exp(-delta / math.Max(temperature, 1e-6))
			if rng.Float64() < threshold {
				best = neighbor
				accept = true
			}
		}

		if accept {
			temperature *= 0.98
		} else {
			temperature *= 0.995
		}
		if temperature < 1e-3 {
			temperature = 1e-3
		}
		best.temperature = temperature
	}

	return best, improvements
}

func (s *solution) destroy(count int, rng *rand.Rand) []int {
	if count <= 0 {
		return []int{}
	}

	removed := make([]int, 0, count)
	maxAttempts := count * len(s.routes) * 2
	attempts := 0

	for len(removed) < count {
		if attempts > maxAttempts {
			break
		}
		attempts++
		shopperIdx := rng.Intn(len(s.routes))
		if len(s.routes[shopperIdx]) == 0 {
			continue
		}
		orderPos := rng.Intn(len(s.routes[shopperIdx]))
		orderID := s.routes[shopperIdx][orderPos]
		s.routes[shopperIdx] = append(s.routes[shopperIdx][:orderPos], s.routes[shopperIdx][orderPos+1:]...)
		removed = append(removed, orderID)
	}

	return removed
}

func (s *solution) repair(removed []int, cache *distanceCache, opts normalizedOptions, rng *rand.Rand) {
	if len(removed) == 0 {
		return
	}

	for _, orderIdx := range removed {
		type insertionOption struct {
			shopper int
			pos     int
			delta   float64
		}

		options := make([]insertionOption, 0, len(cache.shoppers)*2)

		for shopperIdx := range s.routes {
			route := s.routes[shopperIdx]
			if !cache.hasCapacity(shopperIdx, len(route)) {
				continue
			}
			if len(route) == 0 {
				delta := cache.shopperToOrder[shopperIdx][orderIdx]
				options = append(options, insertionOption{shopper: shopperIdx, pos: 0, delta: delta})
				continue
			}
			for pos := 0; pos <= len(route); pos++ {
				delta := insertionDelta(cache, shopperIdx, route, orderIdx, pos)
				options = append(options, insertionOption{
					shopper: shopperIdx,
					pos:     pos,
					delta:   delta,
				})
			}
		}

		if len(options) == 0 {
			fallbackShopper := rng.Intn(len(s.routes))
			route := s.routes[fallbackShopper]
			route = append(route, orderIdx)
			s.routes[fallbackShopper] = route
			continue
		}

		sort.Slice(options, func(i, j int) bool { return options[i].delta < options[j].delta })
		rcl := opts.rclSize
		if rcl < 1 {
			rcl = 1
		}
		if rcl > len(options) {
			rcl = len(options)
		}
		choice := options[rng.Intn(rcl)]
		route := s.routes[choice.shopper]
		if choice.pos >= len(route) {
			route = append(route, orderIdx)
		} else {
			route = append(route[:choice.pos], append([]int{orderIdx}, route[choice.pos:]...)...)
		}
		s.routes[choice.shopper] = route
	}
}

func insertionDelta(cache *distanceCache, shopperIdx int, route []int, orderIdx int, pos int) float64 {
	prevToOrder := 0.0
	if pos == 0 {
		prevToOrder = cache.shopperToOrder[shopperIdx][orderIdx]
	} else {
		prevToOrder = cache.orderToOrder[route[pos-1]][orderIdx]
	}

	orderToNext := 0.0
	if pos == len(route) {
		orderToNext = 0
	} else {
		orderToNext = cache.orderToOrder[orderIdx][route[pos]]
	}

	previousToNext := 0.0
	if len(route) == 0 {
		previousToNext = 0
	} else if pos == 0 {
		previousToNext = cache.shopperToOrder[shopperIdx][route[0]]
	} else if pos == len(route) {
		previousToNext = 0
	} else {
		previousToNext = cache.orderToOrder[route[pos-1]][route[pos]]
	}

	return prevToOrder + orderToNext - previousToNext
}

type distanceCache struct {
	shopperToOrder  [][]float64
	orderToOrder    [][]float64
	orders          []models.Order
	shoppers        []models.Shopper
	totalOrders     int
	capacities      []int
	randomReference float64
}

func newDistanceCache(orders []models.Order, shoppers []models.Shopper) *distanceCache {
	orderCount := len(orders)
	shopperCount := len(shoppers)

	shopperToOrder := make([][]float64, shopperCount)
	for i := range shopperToOrder {
		shopperToOrder[i] = make([]float64, orderCount)
		for j := range orders {
			shopperToOrder[i][j] = optimizer.HaversineDistance(
				shoppers[i].Lat, shoppers[i].Lng,
				orders[j].Lat, orders[j].Lng,
			)
		}
	}

	orderToOrder := make([][]float64, orderCount)
	for i := range orderToOrder {
		orderToOrder[i] = make([]float64, orderCount)
		for j := range orders {
			if i == j {
				orderToOrder[i][j] = 0
			} else {
				orderToOrder[i][j] = optimizer.HaversineDistance(
					orders[i].Lat, orders[i].Lng,
					orders[j].Lat, orders[j].Lng,
				)
			}
		}
	}

	capacities := make([]int, shopperCount)
	for i, shopper := range shoppers {
		if shopper.Capacity <= 0 {
			capacities[i] = -1
		} else {
			capacities[i] = shopper.Capacity
		}
	}

	randomReference := computeBaselineDistance(orders, shoppers)

	return &distanceCache{
		shopperToOrder:  shopperToOrder,
		orderToOrder:    orderToOrder,
		orders:          orders,
		shoppers:        shoppers,
		totalOrders:     orderCount,
		capacities:      capacities,
		randomReference: randomReference,
	}
}

func (dc *distanceCache) routeDistance(shopperIdx int, route []int) float64 {
	if len(route) == 0 {
		return 0
	}
	total := dc.shopperToOrder[shopperIdx][route[0]]
	for i := 0; i < len(route)-1; i++ {
		total += dc.orderToOrder[route[i]][route[i+1]]
	}
	return total
}

func (dc *distanceCache) hasCapacity(shopperIdx int, currentLoad int) bool {
	capacity := dc.capacities[shopperIdx]
	return capacity < 0 || currentLoad < capacity
}

func computeBaselineDistance(orders []models.Order, shoppers []models.Shopper) float64 {
	if len(orders) == 0 || len(shoppers) == 0 {
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
		totalDist += optimizer.HaversineDistance(shopper.Lat, shopper.Lng, order.Lat, order.Lng)

		if i > 0 && (i-1)/ordersPerShopper == shopperIdx {
			prevOrder := orders[i-1]
			totalDist += optimizer.HaversineDistance(prevOrder.Lat, prevOrder.Lng, order.Lat, order.Lng)
		}
	}

	return totalDist
}
