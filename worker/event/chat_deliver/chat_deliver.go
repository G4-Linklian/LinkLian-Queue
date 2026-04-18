package chatdeliver

import (
	"worker/models"
	"worker/utils"
	"worker/utils/logger"
)

func Handle(event models.SocketEvent) error {
	logger.Log("Worker: Emitting CHAT_DELIVER to socket server", "ChatDeliver")

	chatId := event.Payload["chat_id"]
	senderId := event.Payload["sender_id"]
	replyId := event.Payload["reply_id"]
	content, _ := event.Payload["content"].(string)
	fileList := utils.GetStringArray(event.Payload, "file_url")

	err := utils.EmitToSocket("CHAT_DELIVER", map[string]interface{}{
		"chat_id":         chatId,
		"sender_id":       senderId,
		"sender_name":     event.Payload["sender_name"],
		"receive_user_id": event.Payload["receive_user_id"],
		"notification_id": event.Payload["notification_id"],
		"content":         content,
		"reply_id":        replyId,
		"file_url":        fileList,
		"created_at":      event.Payload["created_at"],
	}, "/ws/chat")
	if err != nil {
		logger.Error("Emit socket error", "ChatDeliver", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("State: CHAT_DELIVER Processed Successfully", "ChatDeliver")
	return nil
}
