package notification

import (
	"worker/models"
	"worker/utils"
	"worker/utils/logger"
)

func HandleNotification(event models.SocketEvent) error {
	logger.Log("Worker: Emitting SEND_NOTIFICATION to socket server", "Notification")

	chatId := event.Payload["chat_id"]
	senderId := event.Payload["sender_id"]
	body, _ := event.Payload["body"].(string)
	title, _ := event.Payload["title"].(string)
	receiverIdList := utils.GetStringArray(event.Payload, "target_user_sys_ids")

	err := utils.EmitToSocket("SEND_NOTIFICATION", map[string]interface{}{
		"chat_id":    chatId,
		"sender_id":  senderId,
		"title":    title,
		"body":    body,
		"target_user_sys_ids": receiverIdList,
		"created_at": event.Payload["created_at"],
	}, "/ws/noti")
	if err != nil {
		logger.Error("Emit socket error", "Notification", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("State: SEND_NOTIFICATION Processed Successfully", "Notification")
	return nil
}
