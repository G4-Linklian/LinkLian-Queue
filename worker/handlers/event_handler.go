package handlers

import (
	"database/sql"
	// "fmt"
	// "github.com/gorilla/websocket"
	"log"
	"worker/models"
	"worker/utils"
	// "github.com/lib/pq"
)

type EventHandler struct {
	DB *sql.DB
}

func NewEventHandler(db *sql.DB) *EventHandler {
	return &EventHandler{DB: db}
}

func (h *EventHandler) ProcessEvent(event models.SocketEvent) error {
	utils.LogInfof("Processing Event Type: %s", event.Type)

	switch event.Type {
	case "CHAT_DELIVER":

		utils.LogInfo("Worker: Emitting CHAT_DELIVER to socket server")

		chatId := event.Payload["chat_id"]
		senderId := event.Payload["sender_id"]
		replyId := event.Payload["reply_id"]
		content, _ := event.Payload["content"].(string)
		fileList := utils.GetStringArray(event.Payload, "file_url")

		err := utils.EmitToSocket("CHAT_DELIVER", map[string]interface{}{
			"chat_id":    chatId,
			"sender_id":  senderId,
			"content":    content,
			"reply_id":   replyId,
			"file_url":   fileList,
			"created_at": event.Payload["created_at"],
		})

		if err != nil {
			utils.LogErrorf("Emit socket error: %v", err)
			return err
		}

		utils.LogInfo("State: CHAT_SEND Processed Successfully")
		return nil

	case "NOTIFY_ALERT":
		utils.LogInfo("State: Processing NOTIFY_ALERT")
		log.Println(event.Payload)
		utils.LogInfo("Inserted system_alert")

	default:
		utils.LogErrorf("State: Unknown Event Type: %s", event.Type)
	}
	return nil
}
