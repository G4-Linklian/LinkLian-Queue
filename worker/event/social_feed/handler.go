package socialfeed

import (
	"worker/models"
	"worker/utils"
	"worker/utils/logger"
)

// Handle รับ NOTIFICATION event ที่มี ref_type = "feed-post"
func Handle(event models.SocketEvent) error {
	logger.Log("Emitting social-feed notification", "SocialFeedHandler")

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
		"section_id":       payload["section_id"],
	}, "/ws/notification")
	if err != nil {
		logger.Error("Emit socket error", "SocialFeedHandler", map[string]interface{}{"error": err})
		return err
	}

	logger.Log("Social-feed notification processed", "SocialFeedHandler")
	return nil
}
