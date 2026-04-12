package qna

import (
	"worker/models"
	"worker/utils"
	"worker/utils/logger"
)

// Handle รับ NOTIFICATION event ที่มี ref_type = "qna-live" หรือ "qna-question"
func Handle(event models.SocketEvent) error {
	logger.Log("Emitting qna notification", "QnaHandler")

	payload := event.Payload

	err := utils.EmitToSocket("NOTIFICATION", map[string]interface{}{
		"notification_id":  payload["notification_id"],
		"receive_user_id":  payload["receive_user_id"],
		"actor_id":         payload["actor_id"],
		"actor_name":       payload["actor_name"],
		"title":            payload["title"],
		"body":             payload["body"],
		"ref_id":           payload["ref_id"],
		"ref_type":         payload["ref_type"],
		"feature":          payload["feature"],
	}, "/ws/notification")
	if err != nil {
		logger.Error("Emit socket error", "QnaHandler", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("QnA notification processed", "QnaHandler")
	return nil
}
