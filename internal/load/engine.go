package load

import (
	"sort"
	"sync"
	"time"

	"github.com/williamug/impactbench/internal/models"
	"github.com/williamug/impactbench/internal/runner"
)

type LoadEngine struct {
	adapter runner.Adapter
}

func NewLoadEngine(adapter runner.Adapter) *LoadEngine {
	return &LoadEngine{adapter: adapter}
}

func (e *LoadEngine) Run(target string, users int, duration time.Duration) (models.Metrics, error) {
	results := make(chan models.Metrics, 1000)
	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Start workers
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					metrics, err := e.adapter.Benchmark(target)
					if err == nil {
						results <- metrics
					} else {
						results <- models.Metrics{
							Errors: models.ErrorMetrics{ErrorRatePercent: 100.0},
						}
					}
				}
			}
		}()
	}

	// Run for duration
	time.Sleep(duration)
	close(stop)
	wg.Wait()
	close(results)

	return e.aggregate(results, duration), nil
}

func (e *LoadEngine) aggregate(results <-chan models.Metrics, duration time.Duration) models.Metrics {
	var totalResponseTime int64
	var totalRequests int64
	var errorCount int64
	var latencies []int64

	for m := range results {
		totalRequests++
		lat := m.ResponseTime.AvgMs
		totalResponseTime += lat
		latencies = append(latencies, lat)
		if m.Errors.ErrorRatePercent > 0 {
			errorCount++
		}
	}

	if totalRequests == 0 {
		return models.Metrics{}
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	avg := totalResponseTime / totalRequests
	p95 := latencies[int(float64(len(latencies))*0.95)]
	min := latencies[0]
	max := latencies[len(latencies)-1]

	return models.Metrics{
		ResponseTime: models.ResponseTimeMetrics{
			AvgMs: avg,
			MinMs: min,
			MaxMs: max,
			P95Ms: p95,
		},
		Throughput: models.ThroughputMetrics{
			RequestsPerSecond: float64(totalRequests) / duration.Seconds(),
		},
		Errors: models.ErrorMetrics{
			ErrorRatePercent: (float64(errorCount) / float64(totalRequests)) * 100.0,
		},
	}
}
