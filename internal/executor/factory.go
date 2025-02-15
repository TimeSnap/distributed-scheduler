package executor

import (
	"fmt"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
)

type Factory interface {
	NewExecutor(job *model.Job, options ...Option) (Executor, error)
}

type factory struct {
	client HttpClient
}

func NewFactory(client HttpClient) Factory {
	return &factory{
		client: client,
	}
}

// Option is a function that modifies an executor before it is returned (e.g. WithRetry)
type Option func(executor Executor) Executor

func (f *factory) NewExecutor(job *model.Job, options ...Option) (Executor, error) {

	var executor Executor
	switch job.Type {
	case model.JobTypeHTTP:
		executor = &httpExecutor{Client: f.client}
	case model.JobTypeAMQP:
		executor = &amqpExecutor{}
	default:
		return nil, fmt.Errorf("unknown job type: %v", job.Type)
	}

	for _, option := range options {
		executor = option(executor)
	}

	return executor, nil
}
