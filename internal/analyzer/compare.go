package analyzer

import (
	"time"

	"github.com/williamug/impactbench/internal/models"
)

func Compare(baseline, current models.Benchmark) models.Comparison {
	delta := models.ComparisonDelta{
		ResponseTimeAvg: calculateDelta(float64(baseline.Metrics.ResponseTime.AvgMs), float64(current.Metrics.ResponseTime.AvgMs)),
		QueryCount:      calculateDelta(baseline.Metrics.Database.QueryCountAvg, current.Metrics.Database.QueryCountAvg),
		MemoryAvg:       calculateDelta(baseline.Metrics.Memory.AvgMB, current.Metrics.Memory.AvgMB),
		ErrorRate:       calculateDelta(baseline.Metrics.Errors.ErrorRatePercent, current.Metrics.Errors.ErrorRatePercent),
	}

	verdict := "STABLE"
	if delta.ResponseTimeAvg < -5 {
		verdict = "IMPROVED"
	} else if delta.ResponseTimeAvg > 5 {
		verdict = "REGRESSED"
	}

	return models.Comparison{
		BaselineID: baseline.ID,
		CurrentID:  current.ID,
		Delta:      delta,
		Verdict:    verdict,
		CreatedAt:  time.Now(),
	}
}

func calculateDelta(baseline, current float64) float64 {
	if baseline == 0 {
		if current == 0 {
			return 0
		}
		return 100
	}
	return ((current - baseline) / baseline) * 100
}
