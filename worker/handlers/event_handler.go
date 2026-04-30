package handlers

import (
	chatdeliver "worker/event/chat_deliver"
	qadeliver "worker/event/qa_deliver"
	notification "worker/event/notification"
	"worker/models"
	"worker/utils/logger"
)

type EventHandler struct {
	
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

func (h *EventHandler) ProcessEvent(event models.SocketEvent) error {
	logger.Log("Processing Event Type", "EventHandler", map[string]interface{}{"event_type": event.Type})

	switch event.Type {
	case "CHAT_DELIVER":
		logger.Log("Found CHAT_DELIVER case", "EventHandler")
		return chatdeliver.HandleChat(event)

	case "QA_LIVE_STARTED", "QA_LIVE_ENDED":
		return qadeliver.HandleLiveRoom(event)

	case "FILE_CHANGED":
		return qadeliver.HandleFile(event)

	case "QA_NEW_QUESTION", "QA_QUESTION_UPDATED":
		return qadeliver.HandleQuestion(event)

	case "QA_UPVOTED":
		return qadeliver.HandleVote(event)

	case "SEND_NOTIFICATION":
		logger.Log("State: Processing SEND_NOTIFICATION", "EventHandler")
		return notification.HandleNotification(event)

	default:
		logger.Error("State: Unknown Event Type", "EventHandler", map[string]interface{}{"event_type": event.Type})
	}
	return nil
}
