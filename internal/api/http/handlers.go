package http

import (
	"github.com/TimeSnap/distributed-scheduler/internal/service/job"
	"github.com/TimeSnap/distributed-scheduler/internal/store/postgres"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Log     *otelzap.Logger
	DB      *sqlx.DB
	OpenApi OpenApiConfig
}

// Api constructs a http.Handler with all application routes defined.
func Api(router *gin.Engine, cfg APIMuxConfig) {
	// ==================
	// OpenAPI (will only mount if enabled)
	OpenApiRoute(cfg.OpenApi, router)

	// ==================
	// Jobs

	// Create a new PostgresSQL job store
	jobStore := postgres.New(cfg.DB, cfg.Log)

	// Create a new job service with the job store and logger
	jobService := job.NewService(jobStore, cfg.Log)

	// Create a new jobs handler with the job service
	jobsHandler := NewJobsHandler(jobService)

	// Define a group of routes for the jobs endpoint
	JobsRoutesV1(router, jobsHandler)
}
