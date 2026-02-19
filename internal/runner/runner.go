package runner

import (
	"github.com/williamug/impactbench/internal/models"
)

type Adapter interface {
	Benchmark(target string) (models.Metrics, error)
}

type Runner struct {
	adapter Adapter
}

func NewRunner(adapter Adapter) *Runner {
	return &Runner{adapter: adapter}
}

func (r *Runner) Run(target string, label string, project string) (models.Metrics, error) {
	return r.adapter.Benchmark(target)
}
