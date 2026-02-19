package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "impactbench",
	Short: "ImpactBench is a performance benchmarking and regression tool",
	Long:  `A global CLI tool designed to benchmark application performance, perform regression-aware comparisons, and run load tests.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(4) // Runtime failure as per guide
	}
}

func init() {
	// Root flags will go here
}
