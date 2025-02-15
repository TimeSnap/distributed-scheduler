package postgres

import (
	"encoding/json"
	"time"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	"github.com/TimeSnap/distributed-scheduler/internal/pkg/security"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"
)

var encryptor security.Encryptor

func SetEncryptor(e security.Encryptor) {
	encryptor = e
}

type jobDB struct {
	ID           uuid.UUID      `db:"id"`
	Type         string         `db:"type"`
	Status       string         `db:"status"`
	ExecuteAt    null.Time      `db:"execute_at"`
	CronSchedule null.String    `db:"cron_schedule"`
	HTTPJob      []byte         `db:"http_job"`
	AMQPJob      []byte         `db:"amqp_job"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
	NextRun      null.Time      `db:"next_run"`
	LockedUntil  null.Time      `db:"locked_until"`
	LockedBy     null.String    `db:"locked_by"`
	Tags         pq.StringArray `db:"tags"`
}

func toJobDB(j *model.Job) (*jobDB, error) {
	dbJ := &jobDB{
		ID:           j.ID,
		Type:         string(j.Type),
		Status:       string(j.Status),
		ExecuteAt:    j.ExecuteAt,
		CronSchedule: j.CronSchedule,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
		NextRun:      j.NextRun,
		Tags:         j.Tags,
	}

	if j.HTTPJob != nil {
		switch j.HTTPJob.Auth.Type {
		case model.AuthTypeBasic:
			// Encrypt both the username and password before storing them
			encryptedUsername, err := encryptor.Encrypt(j.HTTPJob.Auth.Username.ValueOrZero())
			if err != nil {
				return nil, err
			}
			j.HTTPJob.Auth.Username = null.StringFrom(*encryptedUsername)

			encryptedPassword, err := encryptor.Encrypt(j.HTTPJob.Auth.Password.ValueOrZero())
			if err != nil {
				return nil, err
			}

			j.HTTPJob.Auth.Password = null.StringFrom(*encryptedPassword)
		case model.AuthTypeBearer:
			encryptedToken, err := encryptor.Encrypt(j.HTTPJob.Auth.BearerToken.ValueOrZero())
			if err != nil {
				return nil, err
			}

			j.HTTPJob.Auth.BearerToken = null.StringFrom(*encryptedToken)
		}

		httpJob, err := json.Marshal(j.HTTPJob)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal http job")
		}

		dbJ.HTTPJob = httpJob
	}

	if j.AMQPJob != nil {

		// Encrypt the connection url before storing it as it contains login credentials
		encryptedConnectionUrl, err := encryptor.Encrypt(j.AMQPJob.Connection)
		if err != nil {
			return nil, err
		}

		j.AMQPJob.Connection = *encryptedConnectionUrl

		amqpJob, err := json.Marshal(j.AMQPJob)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal amqp job")
		}

		dbJ.AMQPJob = amqpJob
	}

	return dbJ, nil
}

func (j *jobDB) ToJob() (*model.Job, error) {
	job := &model.Job{
		ID:           j.ID,
		Type:         model.JobType(j.Type),
		Status:       model.JobStatus(j.Status),
		ExecuteAt:    j.ExecuteAt,
		CronSchedule: j.CronSchedule,
		CreatedAt:    j.CreatedAt,
		UpdatedAt:    j.UpdatedAt,
		NextRun:      j.NextRun,
		Tags:         j.Tags,
	}

	if err := unmarshalNullableJSON(j.HTTPJob, &job.HTTPJob); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal http job")
	}

	if job.HTTPJob != nil {
		switch job.HTTPJob.Auth.Type {
		case model.AuthTypeBasic:
			// Encrypt both the username and password before storing them
			decryptedUsername, err := encryptor.Decrypt(job.HTTPJob.Auth.Username.ValueOrZero())
			if err != nil {
				return nil, err
			}
			job.HTTPJob.Auth.Username = null.StringFrom(*decryptedUsername)

			decryptedPassword, err := encryptor.Decrypt(job.HTTPJob.Auth.Password.ValueOrZero())
			if err != nil {
				return nil, err
			}
			job.HTTPJob.Auth.Password = null.StringFrom(*decryptedPassword)
		case model.AuthTypeBearer:
			decryptedToken, err := encryptor.Decrypt(job.HTTPJob.Auth.BearerToken.ValueOrZero())
			if err != nil {
				return nil, err
			}

			job.HTTPJob.Auth.BearerToken = null.StringFrom(*decryptedToken)
		}
	}

	if err := unmarshalNullableJSON(j.AMQPJob, &job.AMQPJob); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal amqp job")
	}

	if job.AMQPJob != nil {
		decryptedConnectionUrl, err := encryptor.Decrypt(job.AMQPJob.Connection)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decrypt amqp connection url")
		}

		job.AMQPJob.Connection = *decryptedConnectionUrl
	}

	return job, nil
}

func unmarshalNullableJSON(data []byte, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}

type executionDB struct {
	ID           int         `db:"id"`
	JobID        uuid.UUID   `db:"job_id"`
	Status       string      `db:"status"`
	StartTime    time.Time   `db:"start_time"`
	EndTime      time.Time   `db:"end_time"`
	ErrorMessage null.String `db:"error_message"`
	CreatedAt    time.Time   `db:"created_at"`
}

func (e *executionDB) ToModel() *model.JobExecution {
	return &model.JobExecution{
		ID:           e.ID,
		JobID:        e.JobID,
		Success:      e.Status == string(model.JobExecutionStatusSuccessful),
		StartTime:    e.StartTime,
		EndTime:      e.EndTime,
		ErrorMessage: e.ErrorMessage,
	}
}
