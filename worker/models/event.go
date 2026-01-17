package models

type SocketEvent struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}
