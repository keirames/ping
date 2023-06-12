package ws

import (
	"chatroom/logger"
	"encoding/json"
)

// { type : "chat-room/send-message", payload: { roomId: "", text: "" } }

type dataType struct {
	Type string `json:"type"`
}

type sendMessagePayload struct {
	RoomID string `json:"roomId"`
	Text   string `json:"text"`
}

func eventsHandler(data []byte) error {
	var dt dataType
	err := json.Unmarshal(data, &dt)
	if err != nil {
		return err
	}

	if dt.Type == "chat-room/send-message" {
		var smp sendMessagePayload
		if err := json.Unmarshal(data, &smp); err != nil {
			logger.L.Error().Err(err).Msg("data from events is invalid")
			return err
		}
	}

	return nil
}
