package executor

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/TimeSnap/distributed-scheduler/internal/model"
	error2 "github.com/TimeSnap/distributed-scheduler/internal/pkg/error"
	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpExecutor struct{}

func (ae *amqpExecutor) Execute(ctx context.Context, j *model.Job) error {
	// Create a new AMQP connection
	conn, err := amqp.Dial(j.AMQPJob.Connection)
	if err != nil {
		return fmt.Errorf("failed to connect to AMQP: %w", err)
	}
	defer conn.Close()

	// Create a new AMQP channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	var body []byte

	if j.AMQPJob.BodyEncoding != nil {
		switch *j.AMQPJob.BodyEncoding {
		case model.BodyEncodingBase64:
			body, err = base64.StdEncoding.DecodeString(j.AMQPJob.Body)
			if err != nil {
				return fmt.Errorf("failed to decode body: %w", err)
			}
		default:
			return error2.ErrInvalidBodyEncoding
		}

	} else {
		body = []byte(j.AMQPJob.Body)
	}

	// Publish a message to the exchange
	err = ch.PublishWithContext(
		ctx,
		j.AMQPJob.Exchange,   // exchange
		j.AMQPJob.RoutingKey, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: j.AMQPJob.ContentType,
			Headers:     j.AMQPJob.Headers,
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
