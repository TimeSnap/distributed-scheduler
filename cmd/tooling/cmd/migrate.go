package cmd

import (
	"context"
	"time"

	"github.com/TimeSnap/distributed-scheduler/internal/pkg/database"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/database/dbmigrate"
	"github.com/spf13/cobra"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate db to latest version.",
	Run:   migrateRun,
}

var dbConfig database.Config

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVar(&dbConfig.User, "user", "scheduler", "database user")
	migrateCmd.Flags().StringVar(&dbConfig.Password, "pass", "scheduler", "database password")
	migrateCmd.Flags().StringVar(&dbConfig.Host, "host", "localhost:5432", "database host")
	migrateCmd.Flags().StringVar(&dbConfig.Name, "name", "scheduler", "database name")
	migrateCmd.Flags().BoolVar(&dbConfig.DisableTLS, "disable_tls", true, "database sslmode disabled")
	migrateCmd.Flags().IntVar(&dbConfig.MaxIdleConns, "max_idle_conns", 3, "database max idle connections")
	migrateCmd.Flags().IntVar(&dbConfig.MaxOpenConns, "max_open_conns", 2, "database max open connections")
}

func migrateRun(cmd *cobra.Command, args []string) {
	logger := otelzap.L().Sugar()
	db, err := database.Open(dbConfig)
	if err != nil {
		logger.Fatalf("unable to create database connection: %v", err)
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := dbmigrate.Migrate(ctx, db); err != nil {
		logger.Fatalf("unable to migrate the database: %v", err)
		return
	}

	logger.Info("Database migrations complete!")
}
