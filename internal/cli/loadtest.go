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
	"github.com/williamug/impactbench/internal/load"
	"github.com/williamug/impactbench/internal/models"
	"github.com/williamug/impactbench/internal/storage"
)

var loadUsers int
var loadDuration int
var loadUrl string

var loadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "Run a load test",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(2)
		}

		target := loadUrl
		if target == "" {
			target = cfg.BaseURL
		}
		if target == "" {
			fmt.Println("Error: No URL provided.")
			os.Exit(2)
		}

		users := loadUsers
		if users == 0 {
			users = cfg.LoadTest.DefaultUsers
		}
		duration := time.Duration(loadDuration) * time.Second
		if loadDuration == 0 {
			duration = time.Duration(cfg.LoadTest.DefaultDuration) * time.Second
		}

		adapter := http.NewHTTPAdapter()
		engine := load.NewLoadEngine(adapter)

		fmt.Printf("🔥 Starting load test: %d users, %v duration\n", users, duration)
		metrics, err := engine.Run(target, users, duration)
		if err != nil {
			fmt.Printf("Load test failed: %v\n", err)
			os.Exit(4)
		}

		benchmark := models.Benchmark{
			Project:   "default",
			Timestamp: time.Now(),
			Target:    models.Target{Type: "endpoint", Value: target},
			Metrics:   metrics,
			Environment: models.Environment{
				OS:       runtime.GOOS,
				CPUCores: runtime.NumCPU(),
			},
		}

		store, err := storage.NewJSONStore(cfg.Storage.Path + ".json")
		var id string
		if err == nil {
			id, _ = store.SaveBenchmark(benchmark)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Metric", "Value"})
		t.AppendRows([]table.Row{
			{"Snapshot ID", id},
			{"Avg Response Time", fmt.Sprintf("%dms", metrics.ResponseTime.AvgMs)},
			{"P95 Response Time", fmt.Sprintf("%dms", metrics.ResponseTime.P95Ms)},
			{"Throughput", fmt.Sprintf("%.2f RPS", metrics.Throughput.RequestsPerSecond)},
			{"Error Rate", fmt.Sprintf("%.2f%%", metrics.Errors.ErrorRatePercent)},
		})
		t.Render()
	},
}

func init() {
	loadtestCmd.Flags().StringVarP(&loadUrl, "url", "u", "", "URL to test")
	loadtestCmd.Flags().IntVar(&loadUsers, "users", 0, "Number of concurrent users")
	loadtestCmd.Flags().IntVar(&loadDuration, "duration", 0, "Duration in seconds")
	rootCmd.AddCommand(loadtestCmd)
}
