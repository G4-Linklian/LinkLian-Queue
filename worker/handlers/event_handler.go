package handlers

import (
	"fmt"
	chatdeliver "worker/event/chat_deliver"
	"worker/event/community"
	socialfeed "worker/event/social_feed"
	"worker/models"
	"worker/utils/logger"
)

type HandlerFunc func(models.SocketEvent) error

type EventHandler struct{}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// notificationHandlers route NOTIFICATION event ตาม ref_type
// เพิ่ม ref_type ใหม่ได้โดยเพิ่ม entry เดียว ไม่ต้องแตะ switch
var notificationHandlers = map[string]HandlerFunc{
	"feed-post":      socialfeed.Handle,
	"community-post": community.Handle,
	"community":      community.Handle,
}

// eventHandlers route event ตาม type หลัก
var eventHandlers = map[string]HandlerFunc{
	"CHAT_DELIVER": chatdeliver.Handle,
}

func (h *EventHandler) ProcessEvent(event models.SocketEvent) error {
	logger.Log("Processing event", "EventHandler", map[string]any{"type": event.Type})

	// NOTIFICATION → route ต่อด้วย ref_type
	if event.Type == "NOTIFICATION" {
		return h.processNotification(event)
	}

	// event type อื่น → หาใน map โดยตรง
	handler, ok := eventHandlers[event.Type]
	if !ok {
		logger.Error("Unknown event type", "EventHandler", map[string]any{"type": event.Type})
		return fmt.Errorf("unknown event type: %s", event.Type)
	}

	return handler(event)
}

func (h *EventHandler) processNotification(event models.SocketEvent) error {
	refType, _ := event.Payload["ref_type"].(string)

	handler, ok := notificationHandlers[refType]
	if !ok {
		logger.Error("Unknown ref_type", "EventHandler", map[string]any{"ref_type": refType})
		return fmt.Errorf("unknown notification ref_type: %s", refType)
	}

	return handler(event)
}
