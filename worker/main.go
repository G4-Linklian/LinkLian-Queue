package main

import (
	"time"
	// "worker/database"
	"worker/handlers"
	"worker/rabbitmq"
	// "worker/utils"
	"worker/utils/logger"
)

func main() {
	time.Sleep(5 * time.Second)

	// db := database.InitDB()
	// defer db.Close()

	logger.Log("Setting up Event Handler", "Main")
	handler := handlers.NewEventHandler()

	rabbitmq.InitConsumer(handler)
}
