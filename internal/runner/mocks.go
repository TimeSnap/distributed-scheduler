package runner

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/TimeSnap/distributed-scheduler/internal/executor"
	"github.com/TimeSnap/distributed-scheduler/internal/model"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/metrics"
	"github.com/google/uuid"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type mockJobService struct {
	sync.Mutex
	Jobs   []*model.Job
	GetErr error
	FinErr error
}

func (m *mockJobService) GetJobsToRun(_ context.Context, _ time.Time, _ time.Time, _ string, _ uint) ([]*model.Job, error) {
	m.Lock()
	defer m.Unlock()
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	// return a copy of the jobs (so we can modify the slice)
	jobs := make([]*model.Job, len(m.Jobs))
	copy(jobs, m.Jobs)

	return jobs, nil
}

func (m *mockJobService) FinishJobExecution(ctx context.Context, job *model.Job, _, _ time.Time, _ error) error {
	m.Lock()
	defer m.Unlock()
	if m.FinErr != nil {
		return m.FinErr
	}
	for i, j := range m.Jobs {
		if j.ID == job.ID {
			m.Jobs = append(m.Jobs[:i], m.Jobs[i+1:]...)
			break
		}
	}
	return nil
}

func createMockJobService(getErr, finErr error) *mockJobService {
	return &mockJobService{
		Jobs:   []*model.Job{{ID: uuid.MustParse("0053c6a4-ba8b-404e-8e3c-e3875800ed40")}, {ID: uuid.MustParse("0053c6a4-ba8b-404e-8e3c-e3275800ed40")}, {ID: uuid.MustParse("0053c6a4-ba8b-404e-8e3c-e3875800ed40")}},
		GetErr: getErr,
		FinErr: finErr,
	}
}

type mockJobExecutor struct {
	err error
}

func (m *mockJobExecutor) Execute(_ context.Context, _ *model.Job) error {
	return m.err
}

type mockExecutorFactory struct {
	executeErr error
	factoryErr error
}

func (m *mockExecutorFactory) NewExecutor(_ *model.Job, _ ...executor.Option) (executor.Executor, error) {
	if m.factoryErr != nil {
		return nil, m.factoryErr
	}
	return &mockJobExecutor{err: m.executeErr}, nil
}

func createRunnerWithMockExecutor(interval time.Duration, maxConcurrentJobs int, getErr, finErr, factoryErr, execErr error) *Runner {

	executorFactory := &mockExecutorFactory{executeErr: execErr, factoryErr: factoryErr}
	jobService := createMockJobService(getErr, finErr)

	zapL, _ := zap.NewDevelopment()
	logger := otelzap.New(zapL)

	return New(Config{
		JobService:      jobService,
		ExecutorFactory: executorFactory,
		Log:             logger,
		InstanceId:      "test",
		JobExecution: JobExecutionSettings{
			Interval:          interval,
			MaxConcurrentJobs: maxConcurrentJobs,
		},
		Metrics: metrics.NewRunnerMetrics(observability.MetricsConfig{Enabled: false}),
	})
}

// Check that all the jobs have been processed
func assertJobsProcessed(t *testing.T, jobService *mockJobService) {
	jobService.Lock()
	defer jobService.Unlock()
	if len(jobService.Jobs) != 0 {
		t.Errorf("Expected all jobs to have been processed, but got %d", len(jobService.Jobs))
	}
}
