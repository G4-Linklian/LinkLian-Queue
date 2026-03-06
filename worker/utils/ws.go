package utils

import (
	"encoding/json"
	"os"
	"worker/utils/logger"

	"github.com/gorilla/websocket"
)

func EmitToSocket(eventType string, payload map[string]interface{}, path string) error {
	logger.Debug("Starting Emit to Socket", "EmitToSocket", map[string]interface{}{"event_type": eventType})
	conn, _, err := websocket.DefaultDialer.Dial(
		os.Getenv("SOCKET_SERVER_URL")+path,
		nil,
	)
	if err != nil {
		logger.Error("WebSocket connection error", "EmitToSocket", map[string]interface{}{"error": err})
		return err
	}
	defer conn.Close()

	event := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	logger.Debug("event data", "EmitToSocket", map[string]interface{}{"event": event})

	data, _ := json.Marshal(event)

	return conn.WriteMessage(websocket.TextMessage, data)
}
