package main

import (
	"time"
	"worker/database"
	"worker/handlers"
	"worker/rabbitmq"
	"worker/utils"
)

func main() {
	// รอให้ Network พร้อม (เผื่อ Container RabbitMQ ยังไม่ขึ้น)
	time.Sleep(5 * time.Second)

	// 1. เชื่อมต่อ Azure DB
	db := database.InitDB()
	defer db.Close()

	// 2. Setup Handler logic
	utils.LogInfo("Setting up Event Handler")
	handler := handlers.NewEventHandler(db)

	// 3. เชื่อมต่อ RabbitMQ และเริ่มรับงาน
	rabbitmq.InitConsumer(handler)
}
