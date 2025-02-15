package model

import (
	"net/url"

	error2 "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
)

type BodyEncoding string

const (
	BodyEncodingBase64 BodyEncoding = "base64"
)

func (be *BodyEncoding) Valid() bool {
	if be == nil {
		return true
	}
	switch *be {
	case BodyEncodingBase64:
		return true
	default:
		return false
	}
}

type AMQPJob struct {
	Connection   string                 `json:"connection"`    // e.g., "amqp://guest:guest@localhost:5672/"
	Exchange     string                 `json:"exchange"`      // e.g., "my_exchange"
	RoutingKey   string                 `json:"routing_key"`   // e.g., "my_routing_key"
	Headers      map[string]interface{} `json:"headers"`       // e.g., {"x-delay": 10000}
	Body         string                 `json:"body"`          // e.g., "Hello, world!"
	BodyEncoding *BodyEncoding          `json:"body_encoding"` // e.g., null, "base64"
	ContentType  string                 `json:"content_type"`  // e.g., "text/plain"
}

// Validate validates an AMQPJob struct.
func (amqpJob *AMQPJob) Validate() error {
	if amqpJob == nil {
		return error2.ErrAMQPJobNotDefined
	}

	_, err := url.Parse(amqpJob.Connection)
	if err != nil {
		return error2.ErrAMQPConnectionInvalid
	}

	if amqpJob.Exchange == "" {
		return error2.ErrEmptyExchange
	}

	if amqpJob.RoutingKey == "" {
		return error2.ErrEmptyRoutingKey
	}

	if !amqpJob.BodyEncoding.Valid() {
		return error2.ErrInvalidBodyEncoding
	}

	return nil
}

// RemoveCredentials removes the credentials from the AMQPJob struct.
func (amqpJob *AMQPJob) RemoveCredentials() {
	connectionUrl, err := url.Parse(amqpJob.Connection)
	if err != nil {
		return
	}

	// Only redact the credentials from the connection URL
	connectionUrl.User = nil
	amqpJob.Connection = connectionUrl.String()
}
