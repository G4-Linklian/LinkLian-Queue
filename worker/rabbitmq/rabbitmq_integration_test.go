//go:build integration
// +build integration

package rabbitmq

import (
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMQ_Publish_And_Consume(t *testing.T) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		t.Fatal("RABBITMQ_URL is required")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		t.Fatalf("amqp.Dial error: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		t.Fatalf("conn.Channel error: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ci_test_queue", // name
		false,           // durable
		true,            // auto-delete
		false,           // exclusive
		false,           // no-wait
		nil,
	)
	if err != nil {
		t.Fatalf("QueueDeclare error: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		t.Fatalf("Consume error: %v", err)
	}

	body := "hello-from-ci"
	if err := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,
		false,
		amqp.Publishing{Body: []byte(body)},
	); err != nil {
		t.Fatalf("Publish error: %v", err)
	}

	select {
	case d := <-msgs:
		if string(d.Body) != body {
			t.Fatalf("got %q want %q", string(d.Body), body)
		}
	case <-time.After(8 * time.Second):
		t.Fatal("timeout waiting for message")
	}
}
