package handlers

import (
	"fmt"
	chatdeliver "worker/event/chat_deliver"
	qadeliver "worker/event/qa_deliver"
	"worker/event/community"
	"worker/event/qna"
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
	"feed-post":         socialfeed.Handle,
	"feed-comment":      socialfeed.Handle,
	"community-post":    community.Handle,
	"community-comment": community.Handle,
	"community":         community.Handle,
	"qna-live":          qna.Handle,
	"qna-question":      qna.Handle,
}

// eventHandlers route event ตาม type หลัก
var eventHandlers = map[string]HandlerFunc{
	"CHAT_DELIVER": chatdeliver.Handle,
}

func (h *EventHandler) ProcessEvent(event models.SocketEvent) error {
	logger.Log("Processing Event Type", "EventHandler", map[string]interface{}{"event_type": event.Type})

	switch event.Type {
	case "CHAT_DELIVER":
		logger.Log("Found CHAT_DELIVER case", "EventHandler")
		return chatdeliver.Handle(event)

	case "QA_LIVE_STARTED", "QA_LIVE_ENDED":
		return qadeliver.HandleLiveRoom(event)

	case "FILE_CHANGED":
		return qadeliver.HandleFile(event)

	case "QA_NEW_QUESTION", "QA_QUESTION_UPDATED":
		return qadeliver.HandleQuestion(event)

	case "QA_UPVOTED":
		return qadeliver.HandleVote(event)

	case "NOTIFICATION":
		logger.Log("State: Processing NOTIFICATION", "EventHandler")
		return h.processNotification(event)

	default:
		logger.Error("State: Unknown Event Type", "EventHandler", map[string]interface{}{"event_type": event.Type})
	}
	return nil
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
