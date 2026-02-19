package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/williamug/impactbench/internal/models"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func createTables(db *sql.DB) error {
	benchmarksTable := `
	CREATE TABLE IF NOT EXISTS benchmarks (
		id TEXT PRIMARY KEY,
		label TEXT,
		project TEXT,
		branch TEXT,
		commit_hash TEXT,
		target_value TEXT,
		timestamp DATETIME,
		metrics_json TEXT,
		env_json TEXT
	);`

	comparisonsTable := `
	CREATE TABLE IF NOT EXISTS comparisons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		baseline_id TEXT,
		current_id TEXT,
		delta_json TEXT,
		verdict TEXT,
		regression_detected BOOLEAN,
		created_at DATETIME
	);`

	if _, err := db.Exec(benchmarksTable); err != nil {
		return err
	}
	if _, err := db.Exec(comparisonsTable); err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) SaveBenchmark(b models.Benchmark) (string, error) {
	metricsJSON, err := json.Marshal(b.Metrics)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metrics: %w", err)
	}
	envJSON, err := json.Marshal(b.Environment)
	if err != nil {
		return "", fmt.Errorf("failed to marshal environment: %w", err)
	}

	if b.ID == "" {
		b.ID = fmt.Sprintf("run_%d", time.Now().UnixNano())
	}

	_, err = s.db.Exec(
		"INSERT INTO benchmarks (id, label, project, branch, commit_hash, target_value, timestamp, metrics_json, env_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		b.ID, b.Label, b.Project, b.Branch, b.CommitHash, b.Target.Value, b.Timestamp, string(metricsJSON), string(envJSON),
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert benchmark: %w", err)
	}

	return b.ID, nil
}

func (s *SQLiteStore) GetBenchmark(id string) (models.Benchmark, error) {
	var b models.Benchmark
	var metricsJSON, envJSON string
	var timestamp time.Time

	err := s.db.QueryRow(
		"SELECT id, label, project, branch, commit_hash, target_value, timestamp, metrics_json, env_json FROM benchmarks WHERE id = ?",
		id,
	).Scan(&b.ID, &b.Label, &b.Project, &b.Branch, &b.CommitHash, &b.Target.Value, &timestamp, &metricsJSON, &envJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return b, fmt.Errorf("benchmark not found: %s", id)
		}
		return b, fmt.Errorf("failed to query benchmark: %w", err)
	}

	b.Timestamp = timestamp
	if err := json.Unmarshal([]byte(metricsJSON), &b.Metrics); err != nil {
		return b, fmt.Errorf("failed to unmarshal metrics: %w", err)
	}
	if err := json.Unmarshal([]byte(envJSON), &b.Environment); err != nil {
		return b, fmt.Errorf("failed to unmarshal environment: %w", err)
	}

	return b, nil
}

func (s *SQLiteStore) GetLatestBenchmarks(limit int) ([]models.Benchmark, error) {
	rows, err := s.db.Query(
		"SELECT id, label, project, branch, commit_hash, target_value, timestamp, metrics_json, env_json FROM benchmarks ORDER BY timestamp DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest benchmarks: %w", err)
	}
	defer rows.Close()

	var benchmarks []models.Benchmark
	for rows.Next() {
		var b models.Benchmark
		var metricsJSON, envJSON string
		var timestamp time.Time
		if err := rows.Scan(&b.ID, &b.Label, &b.Project, &b.Branch, &b.CommitHash, &b.Target.Value, &timestamp, &metricsJSON, &envJSON); err != nil {
			return nil, fmt.Errorf("failed to scan benchmark: %w", err)
		}
		b.Timestamp = timestamp
		if err := json.Unmarshal([]byte(metricsJSON), &b.Metrics); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metrics: %w", err)
		}
		if err := json.Unmarshal([]byte(envJSON), &b.Environment); err != nil {
			return nil, fmt.Errorf("failed to unmarshal environment: %w", err)
		}
		benchmarks = append(benchmarks, b)
	}
	return benchmarks, nil
}

func (s *SQLiteStore) SaveComparison(c models.Comparison) error {
	deltaJSON, err := json.Marshal(c.Delta)
	if err != nil {
		return fmt.Errorf("failed to marshal delta: %w", err)
	}

	_, err = s.db.Exec(
		"INSERT INTO comparisons (baseline_id, current_id, delta_json, verdict, regression_detected, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		c.BaselineID, c.CurrentID, string(deltaJSON), c.Verdict, c.RegressionDetected, c.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert comparison: %w", err)
	}

	return nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
