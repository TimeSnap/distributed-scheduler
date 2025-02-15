package cmd

import (
	"os"

	"github.com/TimeSnap/distributed-scheduler/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "CLI tool for managing the scheduler.",
}

func Execute() {
	cobra.OnInitialize(logger.SetupLogging)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
