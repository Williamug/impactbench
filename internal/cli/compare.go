package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/williamug/impactbench/internal/analyzer"
	"github.com/williamug/impactbench/internal/config"
	"github.com/williamug/impactbench/internal/storage"
)

var baselineID string
var currentID string

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare two benchmarks",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(2)
		}

		store, err := storage.NewJSONStore(cfg.Storage.Path + ".json")
		if err != nil {
			fmt.Printf("Storage error: %v\n", err)
			os.Exit(3)
		}

		baseline, err := store.GetBenchmark(baselineID)
		if err != nil {
			fmt.Printf("Error fetching baseline (%s): %v\n", baselineID, err)
			os.Exit(3)
		}

		current, err := store.GetBenchmark(currentID)
		if err != nil {
			fmt.Printf("Error fetching current (%s): %v\n", currentID, err)
			os.Exit(3)
		}

		comparison := analyzer.Compare(baseline, current)
		err = store.SaveComparison(comparison)
		if err != nil {
			fmt.Printf("Failed to save comparison: %v\n", err)
			os.Exit(3)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Metric", "Baseline", "Current", "Delta"})
		t.AppendRows([]table.Row{
			{"Response Time", fmt.Sprintf("%dms", baseline.Metrics.ResponseTime.AvgMs), fmt.Sprintf("%dms", current.Metrics.ResponseTime.AvgMs), fmt.Sprintf("%.2f%%", comparison.Delta.ResponseTimeAvg)},
			{"Memory Usage", fmt.Sprintf("%.1fMB", baseline.Metrics.Memory.AvgMB), fmt.Sprintf("%.1fMB", current.Metrics.Memory.AvgMB), fmt.Sprintf("%.2f%%", comparison.Delta.MemoryAvg)},
			{"Error Rate", fmt.Sprintf("%.2f%%", baseline.Metrics.Errors.ErrorRatePercent), fmt.Sprintf("%.2f%%", current.Metrics.Errors.ErrorRatePercent), fmt.Sprintf("%.2f%%", comparison.Delta.ErrorRate)},
		})
		t.AppendFooter(table.Row{"Verdict", "", "", comparison.Verdict})
		t.Render()
	},
}

func init() {
	compareCmd.Flags().StringVar(&baselineID, "baseline", "", "Baseline snapshot ID")
	compareCmd.Flags().StringVar(&currentID, "current", "", "Current snapshot ID")
	compareCmd.MarkFlagRequired("baseline")
	compareCmd.MarkFlagRequired("current")
	rootCmd.AddCommand(compareCmd)
}
