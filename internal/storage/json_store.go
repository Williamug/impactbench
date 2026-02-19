package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/williamug/impactbench/internal/models"
)

type JSONStore struct {
	path string
	mu   sync.RWMutex
	data struct {
		Benchmarks  []models.Benchmark  `json:"benchmarks"`
		Comparisons []models.Comparison `json:"comparisons"`
	}
}

func NewJSONStore(path string) (*JSONStore, error) {
	s := &JSONStore{path: path}
	if _, err := os.Stat(path); err == nil {
		file, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read store file: %w", err)
		}
		if err := json.Unmarshal(file, &s.data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal store data: %w", err)
		}
	}
	return s, nil
}

func (s *JSONStore) SaveBenchmark(b models.Benchmark) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If ID is not set, generate one (for simpler usage)
	if b.ID == "" {
		b.ID = fmt.Sprintf("run_%d", len(s.data.Benchmarks)+1)
	}
	s.data.Benchmarks = append(s.data.Benchmarks, b)

	if err := s.save(); err != nil {
		return "", err
	}
	return b.ID, nil
}

func (s *JSONStore) GetBenchmark(id string) (models.Benchmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, b := range s.data.Benchmarks {
		if b.ID == id {
			return b, nil
		}
	}
	return models.Benchmark{}, fmt.Errorf("benchmark not found: %s", id)
}

func (s *JSONStore) GetLatestBenchmarks(limit int) ([]models.Benchmark, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.data.Benchmarks)
	if n == 0 {
		return nil, fmt.Errorf("no benchmarks found")
	}

	if limit > n {
		limit = n
	}

	result := make([]models.Benchmark, limit)
	for i := 0; i < limit; i++ {
		result[i] = s.data.Benchmarks[n-1-i]
	}
	return result, nil
}

func (s *JSONStore) SaveComparison(c models.Comparison) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data.Comparisons = append(s.data.Comparisons, c)
	return s.save()
}

func (s *JSONStore) save() error {
	file, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal store data: %w", err)
	}
	if err := os.WriteFile(s.path, file, 0644); err != nil {
		return fmt.Errorf("failed to write store file: %w", err)
	}
	return nil
}
