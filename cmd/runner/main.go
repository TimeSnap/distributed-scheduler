package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/TimeSnap/distributed-scheduler/internal/executor"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/database"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/logger"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/metrics"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/security"
	"github.com/TimeSnap/distributed-scheduler/internal/runner"
	"github.com/TimeSnap/distributed-scheduler/internal/service/job"
	"github.com/TimeSnap/distributed-scheduler/internal/store/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	devxCfg "github.com/xBlaz3kx/DevX/configuration"
	devxHttp "github.com/xBlaz3kx/DevX/http"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

const serviceName = "runner"

var serviceInfo = observability.ServiceInfo{
	Name:    serviceName,
	Version: "0.1.2",
}

var configFilePath string

type config struct {
	Observability        observability.Config        `mapstructure:"observability" yaml:"observability" json:"observability"`
	Http                 devxHttp.Configuration      `mapstructure:"http" yaml:"http" json:"http"`
	DB                   database.Config             `mapstructure:"db" yaml:"db" json:"db"`
	ID                   string                      `mapstructure:"id" yaml:"id" json:"id,omitempty"`
	JobExecutionSettings runner.JobExecutionSettings `mapstructure:"jobExecutionSettings" yaml:"jobExecutionSettings" json:"jobExecutionSettings"`
}

var rootCmd = &cobra.Command{
	Use:   "runner",
	Short: "Scheduler runner",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		devxCfg.SetDefaults(serviceName)
		devxCfg.SetupEnv(serviceName)

		viper.SetDefault("storage.encryption.key", "ishouldreallybechanged")
		viper.SetDefault("db.disableTls", true)
		viper.SetDefault("db.maxOpenConns", 1)
		viper.SetDefault("db.maxIdleConns", 10)
		viper.SetDefault("observability.logging.level", observability.LogLevelInfo)

		viper.SetDefault("jobExecutionSettings.maxConcurrentJobs", 100)
		viper.SetDefault("jobExecutionSettings.interval", time.Second*10)
		viper.SetDefault("jobExecutionSettings.maxJobLockTime", time.Minute)

		devxCfg.InitConfig(configFilePath, "./config", ".")

		postgres.SetEncryptor(security.NewEncryptorFromEnv())
	},
	Run: runCmd,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "config file (default is $HOME/.config/runner.yaml)")
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

func main() {
	cobra.OnInitialize(logger.SetupLogging)
	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func runCmd(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()

	// Read the configuration
	cfg := &config{}
	devxCfg.GetConfiguration(viper.GetViper(), cfg)

	obs, err := observability.NewObservability(ctx, serviceInfo, cfg.Observability)
	if err != nil {
		otelzap.L().Fatal("failed to initialize observability", zap.Error(err))
	}

	log := otelzap.L()

	// App Starting
	log.Info("Starting the runner", zap.String("version", serviceInfo.Version))
	defer log.Info("shutdown complete")

	log.Info("Using config", zap.Any("config", cfg))

	// Database
	log.Info("Connecting to the database", zap.String("host", cfg.DB.Host))
	db, err := database.Open(database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		log.Fatal("Unable to establish DB connection", zap.Error(err))
	}

	defer func() {
		log.Info("Closing the database connection")
		_ = db.Close()
	}()

	// Start Runner Service
	log.Info("Starting runner service")

	store := postgres.New(db, log)

	jobService := job.NewService(store, log)

	executorFactory := executor.NewFactory(&http.Client{Timeout: 30 * time.Second})

	runner := runner.New(runner.Config{
		JobService:      jobService,
		Metrics:         metrics.NewRunnerMetrics(cfg.Observability.Metrics),
		Log:             log,
		ExecutorFactory: executorFactory,
		InstanceId:      cfg.ID,
		JobExecution:    cfg.JobExecutionSettings,
	})
	runner.Start()

	httpServer := devxHttp.NewServer(cfg.Http, obs)
	databaseCheck := database.NewHealthChecker(db)
	httpServer.Run(databaseCheck)

	//nolint:all
	<-ctx.Done()
	log.Info("Shutting down the runner")

	httpServer.Shutdown()

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// stop the runner
	runner.Stop(ctx)
}
