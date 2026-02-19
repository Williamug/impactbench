package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/williamug/impactbench/internal/analyzer"
	"github.com/williamug/impactbench/internal/config"
	"github.com/williamug/impactbench/internal/regression"
	"github.com/williamug/impactbench/internal/storage"
)

var failOnRegression bool

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Perform automatic performance review",
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

		latest, err := store.GetLatestBenchmarks(2)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(3)
		}

		if len(latest) < 2 {
			fmt.Println("Error: Need at least two benchmarks to perform a review.")
			os.Exit(3)
		}

		current := latest[0]
		baseline := latest[1]

		fmt.Printf("🔍 Reviewing: %s vs %s\n", current.ID, baseline.ID)

		comp := analyzer.Compare(baseline, current)
		thresholds := regression.Thresholds{
			ResponseTime: cfg.Thresholds.ResponseTime,
			ErrorRate:    cfg.Thresholds.ErrorRate,
		}

		result := regression.Evaluate(comp.Delta, thresholds)

		if result.IsRegression {
			fmt.Println("⚠️  REGRESSION DETECTED")
			for _, v := range result.Violations {
				fmt.Printf("- %s\n", v)
			}
			if failOnRegression {
				os.Exit(1)
			}
		} else {
			fmt.Println("✅ Performance within thresholds.")
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Metric", "Delta", "Status"})
		t.AppendRow(table.Row{"Response Time", fmt.Sprintf("%.2f%%", comp.Delta.ResponseTimeAvg), comp.Verdict})
		t.Render()
	},
}

func init() {
	reviewCmd.Flags().BoolVar(&failOnRegression, "fail-on-regression", false, "Exit with code 1 if regression detected")
	rootCmd.AddCommand(reviewCmd)
}
