package test

import (
	"testing"
	"time"

	"github.com/williamug/impactbench/internal/load"
	"github.com/williamug/impactbench/internal/models"
)

type MockAdapter struct {
	latency int64
}

func (m *MockAdapter) Benchmark(target string) (models.Metrics, error) {
	return models.Metrics{
		ResponseTime: models.ResponseTimeMetrics{AvgMs: m.latency},
	}, nil
}

func TestLoadEngine(t *testing.T) {
	adapter := &MockAdapter{latency: 50}
	engine := load.NewLoadEngine(adapter)

	// Short run for testing
	metrics, err := engine.Run("http://test", 2, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	if metrics.ResponseTime.AvgMs != 50 {
		t.Errorf("Expected 50ms avg latency, got %dms", metrics.ResponseTime.AvgMs)
	}

	if metrics.Throughput.RequestsPerSecond <= 0 {
		t.Errorf("Expected positive throughput, got %f", metrics.Throughput.RequestsPerSecond)
	}
}
