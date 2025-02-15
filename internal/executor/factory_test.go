package executor

import (
	"net/http"
	"testing"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutor(t *testing.T) {
	j := &model.Job{
		Type: model.JobTypeHTTP,
		HTTPJob: &model.HTTPJob{
			Method: "GET",
			URL:    "http://www.example.com",
		},
	}

	factory := NewFactory(&http.Client{})

	executor, err := factory.NewExecutor(j)
	assert.Nil(t, err)
	assert.IsType(t, &httpExecutor{}, executor)

	j.Type = model.JobTypeAMQP
	executor, err = factory.NewExecutor(j)
	assert.Nil(t, err)
	assert.IsType(t, &amqpExecutor{}, executor)

	j.Type = "unknown"
	executor, err = factory.NewExecutor(j)
	assert.NotNil(t, err)
	assert.Nil(t, executor)

	j.Type = model.JobTypeHTTP
	executor, err = factory.NewExecutor(j, WithRetry)
	assert.Nil(t, err)
	assert.IsType(t, &retryExecutor{}, executor)
}
