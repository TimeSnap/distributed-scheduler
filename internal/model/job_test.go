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
