package models

import "time"

type ComparisonDelta struct {
	ResponseTimeAvg float64 `json:"response_time_avg"`
	QueryCount      float64 `json:"query_count"`
	MemoryAvg       float64 `json:"memory_avg"`
	ErrorRate       float64 `json:"error_rate"`
}

type Comparison struct {
	BaselineID         string          `json:"baseline"`
	CurrentID          string          `json:"candidate"`
	Delta              ComparisonDelta `json:"delta"`
	Verdict            string          `json:"verdict"` // IMPROVED, REGRESSED, STABLE
	RegressionDetected bool            `json:"regression_detected"`
	CreatedAt          time.Time       `json:"created_at"`
}
