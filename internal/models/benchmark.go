package models

import "time"

type Environment struct {
	OS       string `json:"os"`
	CPUCores int    `json:"cpu_cores"`
	MemoryMB int    `json:"memory_mb"`
}

type Target struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type ResponseTimeMetrics struct {
	AvgMs int64 `json:"avg_ms"`
	MinMs int64 `json:"min_ms"`
	MaxMs int64 `json:"max_ms"`
	P95Ms int64 `json:"p95_ms"`
}

type ThroughputMetrics struct {
	RequestsPerSecond float64 `json:"requests_per_second"`
}

type DatabaseMetrics struct {
	QueryCountAvg  float64 `json:"query_count_avg"`
	QueryTimeAvgMs float64 `json:"query_time_avg_ms"`
}

type MemoryMetrics struct {
	AvgMB  float64 `json:"avg_mb"`
	PeakMB float64 `json:"peak_mb"`
}

type CPUMetrics struct {
	AvgPercent float64 `json:"avg_percent"`
}

type ErrorMetrics struct {
	ErrorRatePercent float64 `json:"error_rate_percent"`
}

type Metrics struct {
	ResponseTime ResponseTimeMetrics `json:"response_time"`
	Throughput   ThroughputMetrics   `json:"throughput"`
	Database     DatabaseMetrics     `json:"database"`
	Memory       MemoryMetrics       `json:"memory"`
	CPU          CPUMetrics          `json:"cpu"`
	Errors       ErrorMetrics        `json:"errors"`
}

type Benchmark struct {
	ID          string      `json:"id"`
	Label       string      `json:"label"`
	Project     string      `json:"project"`
	Branch      string      `json:"branch"`
	CommitHash  string      `json:"commit_hash"`
	Timestamp   time.Time   `json:"timestamp"`
	Environment Environment `json:"environment"`
	Target      Target      `json:"target"`
	Metrics     Metrics     `json:"metrics"`
	LoadProfile struct {
		VirtualUsers    int `json:"virtual_users"`
		DurationSeconds int `json:"duration_seconds"`
	} `json:"load_profile"`
}
