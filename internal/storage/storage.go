package storage

import (
	"github.com/williamug/impactbench/internal/models"
)

type Storage interface {
	SaveBenchmark(b models.Benchmark) (string, error)
	GetBenchmark(id string) (models.Benchmark, error)
	GetLatestBenchmarks(limit int) ([]models.Benchmark, error)
	SaveComparison(c models.Comparison) error
}
