package cli

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/williamug/impactbench/adapters/http"
	"github.com/williamug/impactbench/internal/config"
	"github.com/williamug/impactbench/internal/models"
	"github.com/williamug/impactbench/internal/runner"
	"github.com/williamug/impactbench/internal/storage"
)

var runUrl string
var runLabel string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a single benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(2)
		}

		target := runUrl
		if target == "" {
			target = cfg.BaseURL
		}
		if target == "" {
			fmt.Println("Error: No URL provided. Use --url or set base_url in config.")
			os.Exit(2)
		}

		adapter := http.NewHTTPAdapter()
		runEngine := runner.NewRunner(adapter)

		fmt.Printf("🚀 Running benchmark for: %s\n", target)
		metrics, err := runEngine.Run(target, runLabel, "default")
		if err != nil {
			fmt.Printf("❌ Benchmark failed: %v\n", err)
			os.Exit(4)
		}

		benchmark := models.Benchmark{
			Label:     runLabel,
			Project:   "default",
			Timestamp: time.Now(),
			Target:    models.Target{Type: "endpoint", Value: target},
			Metrics:   metrics,
			Environment: models.Environment{
				OS:       runtime.GOOS,
				CPUCores: runtime.NumCPU(),
				MemoryMB: 0, // Simplified for now
			},
		}

		store, err := storage.NewJSONStore(cfg.Storage.Path + ".json")
		if err != nil {
			fmt.Printf("Storage error: %v\n", err)
			os.Exit(3)
		}

		id, err := store.SaveBenchmark(benchmark)
		if err != nil {
			fmt.Printf("Failed to save result: %v\n", err)
			os.Exit(3)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Metric", "Value"})
		t.AppendRows([]table.Row{
			{"ID", id},
			{"Label", runLabel},
			{"AVG Latency", fmt.Sprintf("%dms", metrics.ResponseTime.AvgMs)},
			{"Throughput", fmt.Sprintf("%.2f RPS", metrics.Throughput.RequestsPerSecond)},
			{"Error Rate", fmt.Sprintf("%.2f%%", metrics.Errors.ErrorRatePercent)},
		})
		t.Render()
	},
}

func init() {
	runCmd.Flags().StringVarP(&runUrl, "url", "u", "", "URL to benchmark")
	runCmd.Flags().StringVarP(&runLabel, "label", "l", "default", "Label for this benchmark")
	rootCmd.AddCommand(runCmd)
}
