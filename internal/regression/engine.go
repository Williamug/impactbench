package regression

import (
	"fmt"

	"github.com/williamug/impactbench/internal/models"
)

type RegressionResult struct {
	IsRegression bool
	Violations   []string
}

type Thresholds struct {
	ResponseTime float64
	Memory       float64
	ErrorRate    float64
}

func Evaluate(delta models.ComparisonDelta, thresholds Thresholds) RegressionResult {
	result := RegressionResult{
		IsRegression: false,
		Violations:   []string{},
	}

	if delta.ResponseTimeAvg > thresholds.ResponseTime {
		result.IsRegression = true
		result.Violations = append(result.Violations, fmt.Sprintf("Response time regression: %.2f%% (threshold: %.2f%%)", delta.ResponseTimeAvg, thresholds.ResponseTime))
	}

	if delta.ErrorRate > thresholds.ErrorRate {
		result.IsRegression = true
		result.Violations = append(result.Violations, fmt.Sprintf("Error rate regression: %.2f%% (threshold: %.2f%%)", delta.ErrorRate, thresholds.ErrorRate))
	}

	if delta.MemoryAvg > thresholds.Memory {
		result.IsRegression = true
		result.Violations = append(result.Violations, fmt.Sprintf("Memory usage regression: %.2f%% (threshold: %.2f%%)", delta.MemoryAvg, thresholds.Memory))
	}

	return result
}
