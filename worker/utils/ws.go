package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"os"
)

func EmitToSocket(eventType string, payload map[string]interface{}) error {
	conn, _, err := websocket.DefaultDialer.Dial(
		os.Getenv("SOCKET_SERVER_URL") + "/ws/internal",
		nil,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	event := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}

	data, _ := json.Marshal(event)
	return conn.WriteMessage(websocket.TextMessage, data)
}