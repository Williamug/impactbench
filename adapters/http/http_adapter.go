package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/williamug/impactbench/internal/models"
)

type HTTPAdapter struct {
	client *http.Client
}

func NewHTTPAdapter() *HTTPAdapter {
	return &HTTPAdapter{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (a *HTTPAdapter) Benchmark(target string) (models.Metrics, error) {
	start := time.Now()
	resp, err := a.client.Get(target)
	duration := time.Since(start)

	if err != nil {
		return models.Metrics{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	metrics := models.Metrics{
		ResponseTime: models.ResponseTimeMetrics{
			AvgMs: duration.Milliseconds(),
			MinMs: duration.Milliseconds(),
			MaxMs: duration.Milliseconds(),
			P95Ms: duration.Milliseconds(),
		},
		Throughput: models.ThroughputMetrics{
			RequestsPerSecond: 1.0,
		},
	}

	if resp.StatusCode >= 400 {
		metrics.Errors.ErrorRatePercent = 100.0
	}

	return metrics, nil
}
