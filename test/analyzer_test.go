package test

import (
	"testing"

	"github.com/williamug/impactbench/internal/analyzer"
	"github.com/williamug/impactbench/internal/models"
	"github.com/williamug/impactbench/internal/regression"
)

func TestDeltaCalculation(t *testing.T) {
	baseline := models.Benchmark{
		Metrics: models.Metrics{
			ResponseTime: models.ResponseTimeMetrics{AvgMs: 100},
			Throughput:   models.ThroughputMetrics{RequestsPerSecond: 50},
		},
	}
	current := models.Benchmark{
		Metrics: models.Metrics{
			ResponseTime: models.ResponseTimeMetrics{AvgMs: 120},
			Throughput:   models.ThroughputMetrics{RequestsPerSecond: 40},
		},
	}

	comp := analyzer.Compare(baseline, current)

	if comp.Delta.ResponseTimeAvg != 20.0 {
		t.Errorf("Expected 20.0%% response time delta, got %f", comp.Delta.ResponseTimeAvg)
	}
}

func TestRegressionDetection(t *testing.T) {
	delta := models.ComparisonDelta{
		ResponseTimeAvg: 15.0,
		ErrorRate:       0.5,
	}
	thresholds := regression.Thresholds{
		ResponseTime: 10.0,
		ErrorRate:    1.0,
	}

	result := regression.Evaluate(delta, thresholds)

	if !result.IsRegression {
		t.Errorf("Expected regression to be detected")
	}
}
