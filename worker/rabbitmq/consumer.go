package rabbitmq

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"worker/handlers"
	"worker/models"
	"worker/utils/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeName = "linklian_events"

// queueBindings กำหนด queue และ routing key pattern ที่ bind กับ exchange
var queueBindings = map[string]string{
	"chat_events":         "chat.*",
	"qa_events":           "qa_live.#",
	"notification_events": "notification.*",
	"user_events":         "user.*",
}

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
		logger.Warn("RabbitMQ connection failed, retrying...", "InitConsumer", map[string]interface{}{"retry": i + 1})
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", "InitConsumer", map[string]interface{}{"error": err})
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Error("Failed to open channel", "InitConsumer", map[string]interface{}{"error": err})
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// 1. Declare topic exchange
	err = ch.ExchangeDeclare(
		exchangeName, "topic", true, false, false, false, nil,
	)
	if err != nil {
		logger.Error("Failed to declare exchange", "InitConsumer", map[string]interface{}{"error": err})
		log.Fatal("Failed to declare exchange:", err)
	}

	// 2. Declare queues และ bind กับ exchange
	for queue, routingKey := range queueBindings {
		_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			logger.Error("Failed to declare queue", "InitConsumer", map[string]interface{}{"queue": queue, "error": err})
			log.Fatalf("Failed to declare queue %s: %v", queue, err)
		}

		err = ch.QueueBind(queue, routingKey, exchangeName, false, nil)
		if err != nil {
			logger.Error("Failed to bind queue", "InitConsumer", map[string]interface{}{"queue": queue, "routing_key": routingKey, "error": err})
			log.Fatalf("Failed to bind queue %s: %v", queue, err)
		}

		logger.Log("Queue bound", "InitConsumer", map[string]interface{}{"queue": queue, "routing_key": routingKey})
	}

	// 3. Fair dispatch — 1 message ต่อ worker
	err = ch.Qos(1, 0, false)
	if err != nil {
		logger.Error("Failed to set QoS", "InitConsumer", map[string]interface{}{"error": err})
		log.Fatal("Failed to set QoS:", err)
	}

	// 4. เปิด goroutine consume แต่ละ queue
	for queue := range queueBindings {
		go consumeQueue(ch, queue, handler)
	}

	logger.Log("Worker is ready to consume messages", "InitConsumer", map[string]interface{}{"exchange": exchangeName})
	select {}
}

func consumeQueue(ch *amqp.Channel, queue string, handler *handlers.EventHandler) {
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		logger.Error("Failed to consume queue", "consumeQueue", map[string]interface{}{"queue": queue, "error": err})
		return
	}

	logger.Log("Consuming queue", "consumeQueue", map[string]interface{}{"queue": queue})

	for d := range msgs {
		var event models.SocketEvent
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			logger.Error("Error decoding JSON", "consumeQueue", map[string]interface{}{"queue": queue, "error": err})
			d.Ack(false) // Ack ทิ้งถ้า JSON ผิด
			continue
		}

		err = handler.ProcessEvent(event)
		if err != nil {
			logger.Error("Error processing event", "consumeQueue", map[string]interface{}{"queue": queue, "error": err})
			// d.Nack(false, true) // เปิดถ้าอยากให้ re-queue
		} else {
			d.Ack(false)
		}
	}
}
