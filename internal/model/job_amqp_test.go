package model

import (
	"testing"

	error2 "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
	"github.com/stretchr/testify/assert"
)

func TestAMQPJobValidate(t *testing.T) {
	tests := []struct {
		name string
		job  AMQPJob
		want error
	}{
		{
			name: "valid job",
			job: AMQPJob{
				Connection: "amqp://guest:guest@localhost:5672/",
				Exchange:   "my_exchange",
				RoutingKey: "my_routing_key",
			},
			want: nil,
		},
		{
			name: "invalid job: empty Exchange",
			job: AMQPJob{
				Connection: "amqp://guest:guest@localhost:5672/",
				Exchange:   "",
				RoutingKey: "my_routing_key",
			},
			want: error2.ErrEmptyExchange,
		},
		{
			name: "invalid job: empty RoutingKey",
			job: AMQPJob{
				Connection: "amqp://guest:guest@localhost:5672/",
				Exchange:   "my_exchange",
				RoutingKey: "",
			},
			want: error2.ErrEmptyRoutingKey,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.job.Validate()
			assert.Equal(t, tc.want, got)
		})
	}
}
