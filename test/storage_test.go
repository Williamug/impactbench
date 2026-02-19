package test

import (
	"os"
	"testing"
	"time"

	"github.com/williamug/impactbench/internal/models"
	"github.com/williamug/impactbench/internal/storage"
)

func TestJSONStore(t *testing.T) {
	tmpFile := "test_store.json"
	defer os.Remove(tmpFile)

	store, err := storage.NewJSONStore(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	benchmark := models.Benchmark{
		ID:        "test_1",
		Project:   "test_proj",
		Timestamp: time.Now(),
		Metrics: models.Metrics{
			ResponseTime: models.ResponseTimeMetrics{AvgMs: 100},
		},
	}

	id, err := store.SaveBenchmark(benchmark)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	if id != "test_1" {
		t.Errorf("Expected ID test_1, got %s", id)
	}

	fetched, err := store.GetBenchmark("test_1")
	if err != nil {
		t.Fatalf("Failed to fetch: %v", err)
	}

	if fetched.Metrics.ResponseTime.AvgMs != 100 {
		t.Errorf("Expected 100ms, got %dms", fetched.Metrics.ResponseTime.AvgMs)
	}
}
