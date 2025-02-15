package executor

import (
	"context"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	"github.com/cenkalti/backoff/v4"
)

// RetryExecutor struct encapsulates an executor and adds retry functionality
type retryExecutor struct {
	executor Executor
}

// WithRetry wraps an executor with a retry mechanism
func WithRetry(executor Executor) Executor {
	return &retryExecutor{executor: executor}
}

const maxRetries = 3

// Execute applies the retry mechanism on the execution of the job
func (re *retryExecutor) Execute(ctx context.Context, job *model.Job) error {
	// Define your backoff strategy
	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), maxRetries)

	// Use the backoff.Retry function with your execute function
	err := backoff.Retry(func() error {
		return re.executor.Execute(ctx, job)
	}, bo)

	return err
}
