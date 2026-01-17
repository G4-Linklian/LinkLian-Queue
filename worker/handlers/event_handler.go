package handlers

import (
	"database/sql"
	// "fmt"
	"log"
	"worker/models"
	"worker/utils"

	"github.com/lib/pq"
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
	case "CHAT_SEND":
		utils.LogInfo("State: Extracting payload data for CHAT_SEND")

		// chatIdPtr := utils.GetIntPointer(event.Payload, "chat_id")
		// senderIdPtr := utils.GetIntPointer(event.Payload, "sender_id")
		// replyIdPtr := utils.GetIntPointer(event.Payload, "reply_id")

		// if chatIdPtr == nil || senderIdPtr == nil {
		// 	utils.LogErrorf("State: Validation Failed - missing chat_id or sender_id. Payload: %v", event.Payload)
		// 	return fmt.Errorf("missing chat_id or sender_id")
		// }

		// content, _ := event.Payload["content"].(string)
		// fileList := utils.GetStringArray(event.Payload, "file_url")

		chatId := event.Payload["chat_id"]
		senderId := event.Payload["sender_id"]
		replyId := event.Payload["reply_id"]
		content, _ := event.Payload["content"].(string)
		fileList := utils.GetStringArray(event.Payload, "file_url")

		utils.LogInfof("State: Preparing to insert message. ChatID: %v, SenderID: %v, ContentLen: %d", chatId, senderId, len(content))

		query := `
			INSERT INTO message (chat_id, sender_id, content, reply_id, file, created_at, flag_valid) 
			VALUES ($1, $2, $3, $4, $5, NOW(), $6)
		`

		utils.LogInfo("State: Executing database insert")
		_, err := h.DB.Exec(query,
			chatId,
			senderId,
			content,
			replyId,
			pq.Array(fileList),
			true,
		)

		if err != nil {
			utils.LogErrorf("State: Database Insert Error: %v", err)
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
