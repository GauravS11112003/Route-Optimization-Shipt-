package models

import "time"

// HybridSolveOptions configures the hybrid solver behaviour.
type HybridSolveOptions struct {
	Iterations            int     `json:"iterations"`
	Workers               int     `json:"workers"`
	CandidatePool         int     `json:"candidatePool"`
	RandomizedListSize    int     `json:"randomizedListSize"`
	DestroyRate           float64 `json:"destroyRate"`
	LocalSearchIterations int     `json:"localSearchIterations"`
	EmitIntervalMillis    int     `json:"emitIntervalMillis"`
	RandomSeed            int64   `json:"randomSeed"`
	UseRealRoutes         bool    `json:"useRealRoutes"`
	ApiKey                string  `json:"apiKey"`
}

// HybridSolveRequest is the request payload used by the hybrid solver endpoint.
type HybridSolveRequest struct {
	Orders   []Order            `json:"orders"`
	Shoppers []Shopper          `json:"shoppers"`
	Options  HybridSolveOptions `json:"options"`
}

// HybridProgress describes an intermediate solver snapshot.
type HybridProgress struct {
	Timestamp           time.Time `json:"timestamp"`
	Iteration           int       `json:"iteration"`
	WorkerID            int       `json:"workerId"`
	BestDistance        float64   `json:"bestDistance"`
	CandidateDistance   float64   `json:"candidateDistance"`
	AcceptedImprovement bool      `json:"acceptedImprovement"`
	ExploredSolutions   int       `json:"exploredSolutions"`
	ImprovementCount    int       `json:"improvementCount"`
	Temperature         float64   `json:"temperature"`
}

// HybridSolverStats captures summary statistics for a solve run.
type HybridSolverStats struct {
	Runtime              time.Duration `json:"runtime"`
	Iterations           int           `json:"iterations"`
	BestIteration        int           `json:"bestIteration"`
	Workers              int           `json:"workers"`
	ExploredSolutions    int           `json:"exploredSolutions"`
	AcceptedImprovements int           `json:"acceptedImprovements"`
}

// HybridSolveResponse is returned when the hybrid solver finishes.
type HybridSolveResponse struct {
	Optimization OptimizeResponse   `json:"optimization"`
	Analytics    *AnalyticsResponse `json:"analytics,omitempty"`
	Stats        HybridSolverStats  `json:"stats"`
	Timeline     []HybridProgress   `json:"timeline"`
}
