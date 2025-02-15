package postgres

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/security"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func init() {
	SetEncryptor(security.NewEncryptor("testkey123456789"))
}

func TestJobDB_ToJob_HTTPJob(t *testing.T) {
	jobDB := &jobDB{
		ID:           uuid.MustParse("a787fa30-2cbe-40de-9a51-f7c9fc43a747"),
		Type:         "http",
		Status:       "scheduled",
		ExecuteAt:    null.TimeFrom(time.Now()),
		CronSchedule: null.StringFrom("0 0 * * *"),
		HTTPJob:      []byte(`{"url": "localhost:3000", "auth": {"type": "none", "password": null, "username": null, "bearer_token": null}, "body": "", "method": "POST", "headers": null}`),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	job, err := jobDB.ToJob()
	require.NoError(t, err)

	assert.Equal(t, job.ID, jobDB.ID)
	assert.Equal(t, job.Type, model.JobType(jobDB.Type))
	assert.Equal(t, job.Status, model.JobStatus(jobDB.Status))
	assert.Equal(t, job.ExecuteAt, jobDB.ExecuteAt)
	assert.Equal(t, job.CronSchedule, jobDB.CronSchedule)
	assert.Nil(t, job.AMQPJob)
	assert.NotNil(t, job.HTTPJob)
	assert.Equal(t, job.CreatedAt, jobDB.CreatedAt)
	assert.Equal(t, job.UpdatedAt, jobDB.UpdatedAt)
	assert.Equal(t, job.NextRun.Valid, false)
}

func TestJobDB_ToJob_AMQPJob(t *testing.T) {

	// AMQP connection must be encrypted
	connection, err := encryptor.Encrypt("amqp://localhost:3000")
	assert.NoError(t, err)
	amqpJob := fmt.Sprintf(`{"connection": "%s", "exchange": "Test", "routing_key": "Test", "headers": {}, "body": "Text Plain", "body_encoding": null, "content_type": "text/plain"}`, *connection)

	jobDB := &jobDB{
		ID:           uuid.MustParse("a787fa30-2cbe-40de-9a51-f7c9fc43a747"),
		Type:         "amqp",
		Status:       "scheduled",
		ExecuteAt:    null.TimeFrom(time.Now()),
		CronSchedule: null.StringFrom("0 0 * * *"),
		AMQPJob:      []byte(amqpJob),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	job, err := jobDB.ToJob()
	require.NoError(t, err)

	assert.Equal(t, job.ID, jobDB.ID)
	assert.Equal(t, job.Type, model.JobType(jobDB.Type))
	assert.Equal(t, job.Status, model.JobStatus(jobDB.Status))
	assert.Equal(t, job.ExecuteAt, jobDB.ExecuteAt)
	assert.Equal(t, job.CronSchedule, jobDB.CronSchedule)
	assert.NotNil(t, job.AMQPJob)
	assert.Nil(t, job.HTTPJob)
	assert.Equal(t, job.CreatedAt, jobDB.CreatedAt)
	assert.Equal(t, job.UpdatedAt, jobDB.UpdatedAt)
	assert.Equal(t, job.NextRun.Valid, false)

	marshalledJob, err := json.Marshal(job.AMQPJob)
	require.NoError(t, err)

	// AMQP connection is decrypted
	assert.JSONEq(t, `{"connection": "amqp://localhost:3000", "exchange": "Test", "routing_key": "Test", "headers": {}, "body": "Text Plain", "body_encoding": null, "content_type": "text/plain"}`, string(marshalledJob))
}
