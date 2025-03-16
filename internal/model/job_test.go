package model

import (
	"testing"
	"time"

	error2 "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestJobTypeValid(t *testing.T) {
	var jobType JobType = "INVALID"

	if jobType.Valid() {
		t.Error("Expected false, got true")
	}

	jobType = JobTypeHTTP

	if !jobType.Valid() {
		t.Error("Expected true, got false")
	}
}

func TestJobStatusValid(t *testing.T) {
	var jobStatus JobStatus = "INVALID"

	if jobStatus.Valid() {
		t.Error("Expected false, got true")
	}

	jobStatus = JobStatusRunning

	if !jobStatus.Valid() {
		t.Error("Expected true, got false")
	}
}

func TestJobValidate(t *testing.T) {
	tests := []struct {
		name string
		job  Job
		want error
	}{
		{
			name: "valid job",
			job: Job{
				ID:        uuid.New(),
				Type:      JobTypeHTTP,
				Status:    JobStatusRunning,
				ExecuteAt: null.TimeFrom(time.Now().Add(time.Minute)),
				HTTPJob: &HTTPJob{
					URL:    "https://example.com",
					Method: "GET",
					Auth: Auth{
						Type: AuthTypeNone,
					},
				},
				CreatedAt: time.Now(),
			},
			want: nil,
		},
		{
			name: "invalid job: missing ID",
			job: Job{
				Type:      JobTypeHTTP,
				Status:    JobStatusRunning,
				ExecuteAt: null.TimeFrom(time.Now().Add(time.Minute)),
				HTTPJob: &HTTPJob{
					URL:    "https://example.com",
					Method: "GET",
					Auth: Auth{
						Type: AuthTypeNone,
					},
				},
				CreatedAt: time.Now(),
			},
			want: error2.ErrInvalidJobID,
		},
		{
			name: "invalid job: http type with nil HTTPJob",
			job: Job{
				ID:        uuid.New(),
				Type:      JobTypeHTTP,
				Status:    JobStatusRunning,
				ExecuteAt: null.TimeFrom(time.Now().Add(time.Minute)),
				CreatedAt: time.Now(),
			},
			want: error2.ErrHTTPJobNotDefined,
		},
		{
			name: "invalid job: unsupported Type",
			job: Job{
				ID:        uuid.New(),
				Type:      "invalid_type",
				Status:    JobStatusRunning,
				ExecuteAt: null.TimeFrom(time.Now().Add(time.Minute)),
				HTTPJob: &HTTPJob{
					URL:    "https://example.com",
					Method: "GET",
					Auth: Auth{
						Type: AuthTypeNone,
					},
				},
				CreatedAt: time.Now(),
			},
			want: error2.ErrInvalidJobType,
		},
		{
			name: "invalid job: invalid cron expression",
			job: Job{
				ID:           uuid.New(),
				Type:         JobTypeHTTP,
				Status:       JobStatusRunning,
				CronSchedule: null.StringFrom("invalid_cron_expression"),
				HTTPJob: &HTTPJob{
					URL:    "https://example.com",
					Method: "GET",
					Auth: Auth{
						Type: AuthTypeNone,
					},
				},
				CreatedAt: time.Now(),
			},
			want: error2.ErrInvalidCronSchedule,
		},
		{
			name: "invalid job: schedule and execute at both defined",
			job: Job{
				ID:           uuid.New(),
				Type:         JobTypeHTTP,
				Status:       JobStatusRunning,
				CronSchedule: null.StringFrom("* * * * *"),
				ExecuteAt:    null.TimeFrom(time.Now().Add(time.Minute)),
				HTTPJob: &HTTPJob{
					URL:    "https://example.com",
					Method: "GET",
					Auth: Auth{
						Type: AuthTypeNone,
					},
				},
				CreatedAt: time.Now(),
			},
			want: error2.ErrInvalidJobSchedule,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.job.Validate()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestRemoveCredentials(t *testing.T) {
	tests := []struct {
		name string
		job  Job
		want Job
	}{
		{
			name: "AMQP job without credentials",
			job: Job{
				Type: JobTypeAMQP,
				AMQPJob: &AMQPJob{
					Connection: "amqp://localhost:5672/",
				},
			},
			want: Job{
				Type: JobTypeAMQP,
				AMQPJob: &AMQPJob{
					Connection: "amqp://localhost:5672/",
				},
			},
		},
		{
			name: "AMQP job with credentials",
			job: Job{
				Type: JobTypeAMQP,
				AMQPJob: &AMQPJob{
					Connection: "amqp://guest:guest@localhost:5672/",
				},
			},
			want: Job{
				Type: JobTypeAMQP,
				AMQPJob: &AMQPJob{
					Connection: "amqp://guest:xxxxx@localhost:5672/",
				},
			},
		},
		{
			name: "HTTP job without any credentials",
			job: Job{
				Type: JobTypeHTTP,
				HTTPJob: &HTTPJob{
					URL:  "https://example.com",
					Auth: Auth{Type: AuthTypeNone},
				},
			},
			want: Job{
				Type: JobTypeHTTP,
				HTTPJob: &HTTPJob{
					URL:  "https://example.com",
					Auth: Auth{Type: AuthTypeNone},
				},
			},
		},
		{
			name: "HTTP job with Bearer token",
			job: Job{
				Type: JobTypeHTTP,
				HTTPJob: &HTTPJob{
					URL: "https://example.com",
					Auth: Auth{
						Type:        AuthTypeBearer,
						BearerToken: null.NewString("imabearertoken123", true),
					},
				},
			},
			want: Job{
				Type: JobTypeHTTP,
				HTTPJob: &HTTPJob{
					URL: "https://example.com",
					Auth: Auth{
						Type:        AuthTypeBearer,
						BearerToken: null.NewString("", false),
					},
				},
			},
		},
		{
			name: "HTTP job with HTTP Basic Auth",
			job: Job{
				Type: JobTypeHTTP,
				HTTPJob: &HTTPJob{
					URL: "https://example.com",
					Auth: Auth{
						Type:     AuthTypeBasic,
						Username: null.NewString("username123", true),
						Password: null.NewString("password123", true),
					},
				},
			},
			want: Job{
				Type: JobTypeHTTP,
				HTTPJob: &HTTPJob{
					URL: "https://example.com",
					Auth: Auth{
						Type:     AuthTypeBasic,
						Username: null.NewString("", false),
						Password: null.NewString("", false)},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			job := tc.job
			job.RemoveCredentials()

			assert.Equal(t, tc.want, job)
		})
	}
}
