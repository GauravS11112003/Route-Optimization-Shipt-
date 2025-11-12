package hybrid

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"shipt-route-optimizer/internal/models"
	"shipt-route-optimizer/internal/optimizer"
)

// Run executes the hybrid GRASP + ALNS solver and emits progress snapshots as they become available.
func Run(
	ctx context.Context,
	orders []models.Order,
	shoppers []models.Shopper,
	req models.HybridSolveOptions,
	emit func(models.HybridProgress),
) (models.HybridSolveResponse, error) {
	opts := normalizeOptions(req)

	if len(orders) == 0 {
		return emptyHybridResponse(shoppers, opts), nil
	}

	if len(shoppers) == 0 {
		return models.HybridSolveResponse{}, errors.New("no shoppers provided")
	}

	dcache := newDistanceCache(orders, shoppers)
	start := time.Now()

	var (
		bestMu            sync.Mutex
		bestSolution      *solution
		bestIteration     int
		bestImprovement   int64
		exploredSolutions atomic.Int64
		acceptedImproves  atomic.Int64
		timelineMu        sync.Mutex
		timeline          []models.HybridProgress
	)

	appendTimeline := func(snapshot models.HybridProgress) {
		timelineMu.Lock()
		defer timelineMu.Unlock()
		timeline = append(timeline, snapshot)
		if opts.candidatePool > 0 && len(timeline) > opts.candidatePool {
			start := len(timeline) - opts.candidatePool
			cp := make([]models.HybridProgress, opts.candidatePool)
			copy(cp, timeline[start:])
			timeline = cp
		}
		if emit != nil {
			emit(snapshot)
		}
	}

	iterationCh := make(chan iterationTask)
	var wg sync.WaitGroup

	for workerID := 0; workerID < opts.workers; workerID++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			rng := rand.New(rand.NewSource(opts.randomSeed + int64(id)*7919))
			lastEmit := time.Time{}

			for task := range iterationCh {
				if ctx.Err() != nil {
					return
				}

				initial := buildInitialSolution(dcache, opts, rng)
				improved, improvementsMade := runLocalSearch(initial, dcache, opts, rng)

				explored := int(exploredSolutions.Add(1))
				if improvementsMade > 0 {
					acceptedImproves.Add(int64(improvementsMade))
				}

				accepted := false
				bestMu.Lock()
				if bestSolution == nil || improved.totalDistance < bestSolution.totalDistance {
					bestSolution = improved.clone()
					bestIteration = task.iteration
					bestImprovement = acceptedImproves.Load()
					accepted = true
				}
				currentBest := bestSolution.totalDistance
				bestMu.Unlock()

				now := time.Now()
				shouldEmit := accepted || lastEmit.IsZero() || now.Sub(lastEmit) >= opts.emitInterval
				if shouldEmit {
					lastEmit = now
					snapshot := models.HybridProgress{
						Timestamp:           now,
						Iteration:           task.iteration,
						WorkerID:            id,
						BestDistance:        math.Round(currentBest*100) / 100,
						CandidateDistance:   math.Round(improved.totalDistance*100) / 100,
						AcceptedImprovement: accepted,
						ExploredSolutions:   explored,
						ImprovementCount:    int(acceptedImproves.Load()),
						Temperature:         improved.temperature,
					}
					appendTimeline(snapshot)
				}
			}
		}(workerID)
	}

	go func() {
		defer close(iterationCh)
		for iter := 0; iter < opts.iterations; iter++ {
			if ctx.Err() != nil {
				return
			}
			iterationCh <- iterationTask{iteration: iter}
		}
	}()

	wg.Wait()

	if ctx.Err() != nil {
		return models.HybridSolveResponse{}, ctx.Err()
	}

	if bestSolution == nil {
		return models.HybridSolveResponse{}, errors.New("solver failed to find a solution")
	}

	assignments := bestSolution.toAssignments(orders, shoppers, dcache)
	optimizer.SortAssignmentsByShopper(assignments)

	analytics := optimizer.AnalyticsFromAssignments(
		orders,
		shoppers,
		assignments,
		opts.useRealRoutes,
		opts.apiKey,
	)

	response := models.HybridSolveResponse{
		Optimization: models.OptimizeResponse{
			Assignments:         assignments,
			TotalDistanceBefore: math.Round(dcache.randomReference*100) / 100,
			TotalDistanceAfter:  math.Round(bestSolution.totalDistance*100) / 100,
		},
		Analytics: analytics,
		Stats: models.HybridSolverStats{
			Runtime:              time.Since(start),
			Iterations:           opts.iterations,
			BestIteration:        bestIteration,
			Workers:              opts.workers,
			ExploredSolutions:    int(exploredSolutions.Load()),
			AcceptedImprovements: int(bestImprovement),
		},
		Timeline: timeline,
	}

	return response, nil
}

type iterationTask struct {
	iteration int
}

type normalizedOptions struct {
	iterations    int
	workers       int
	candidatePool int
	rclSize       int
	destroyRate   float64
	localSearch   int
	emitInterval  time.Duration
	randomSeed    int64
	useRealRoutes bool
	apiKey        string
}

func normalizeOptions(req models.HybridSolveOptions) normalizedOptions {
	opts := normalizedOptions{
		iterations:    req.Iterations,
		workers:       req.Workers,
		candidatePool: req.CandidatePool,
		rclSize:       req.RandomizedListSize,
		destroyRate:   req.DestroyRate,
		localSearch:   req.LocalSearchIterations,
		emitInterval:  time.Duration(req.EmitIntervalMillis) * time.Millisecond,
		randomSeed:    req.RandomSeed,
		useRealRoutes: req.UseRealRoutes,
		apiKey:        req.ApiKey,
	}

	if opts.iterations <= 0 {
		opts.iterations = 400
	}
	if opts.workers <= 0 {
		opts.workers = runtime.NumCPU()
		if opts.workers < 1 {
			opts.workers = 1
		}
	}
	if opts.workers > opts.iterations {
		opts.workers = opts.iterations
		if opts.workers < 1 {
			opts.workers = 1
		}
	}
	if opts.candidatePool <= 0 {
		opts.candidatePool = opts.iterations
	}
	if opts.rclSize <= 0 {
		opts.rclSize = 3
	}
	if opts.destroyRate <= 0 {
		opts.destroyRate = 0.35
	}
	if opts.localSearch <= 0 {
		opts.localSearch = 50
	}
	if opts.emitInterval <= 0 {
		opts.emitInterval = 250 * time.Millisecond
	}
	if opts.randomSeed == 0 {
		opts.randomSeed = time.Now().UnixNano()
	}
	return opts
}

func emptyHybridResponse(shoppers []models.Shopper, opts normalizedOptions) models.HybridSolveResponse {
	return models.HybridSolveResponse{
		Optimization: models.OptimizeResponse{
			Assignments:         []models.Assignment{},
			TotalDistanceBefore: 0,
			TotalDistanceAfter:  0,
		},
		Analytics: nil,
		Stats: models.HybridSolverStats{
			Runtime:              0,
			Iterations:           opts.iterations,
			BestIteration:        0,
			Workers:              opts.workers,
			ExploredSolutions:    0,
			AcceptedImprovements: 0,
		},
		Timeline: []models.HybridProgress{},
	}
}
