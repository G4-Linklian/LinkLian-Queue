package rabbitmq

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"worker/handlers"
	"worker/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitConsumer(handler *handlers.EventHandler) {
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://user:password@rabbitmq:5672/"
	}

	var conn *amqp.Connection
	var err error

	// Loop Retry เผื่อ RabbitMQ ยัง Boot ไม่เสร็จ
	for i := 0; i < 5; i++ {
		conn, err = amqp.Dial(amqpURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ connection failed, retrying... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	qName := "socket_events"

	// ประกาศ Queue (ต้องตรงกับฝั่งคนส่ง)
	_, err = ch.QueueDeclare(
		qName, true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// รับงานทีละ 1 งาน
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatal("Failed to set QoS:", err)
	}

	msgs, err := ch.Consume(
		qName, "", false, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	log.Printf(" [*] Waiting for messages in %s...", qName)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var event models.SocketEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				d.Ack(false) // Ack ทิ้งไปเลยถ้า JSON ผิด
				continue
			}

			err = handler.ProcessEvent(event)
			if err != nil {
				log.Printf("Error inserting to DB: %s", err)
				// d.Nack(false, true) // เปิดบรรทัดนี้ถ้าอยากให้มัน Re-queue กรณี Insert ไม่เข้า
			} else {
				d.Ack(false)
			}
		}
	}()

	<-forever
}
