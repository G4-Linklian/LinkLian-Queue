package handlers

import (
	chatdeliver "worker/event/chat_deliver"
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
		return chatdeliver.Handle(event)

	// case "NOTIFY_ALERT":
	// 	logger.Log("State: Processing NOTIFY_ALERT", "EventHandler")
	// 	log.Println(event.Payload)
	// 	logger.Log("State: Inserted system_alert", "EventHandler")

	default:
		logger.Error("State: Unknown Event Type", "EventHandler", map[string]interface{}{"event_type": event.Type})
	}
	return nil
}
