package ws

import (
	"chatroom/logger"
	"encoding/json"
	"strconv"
)

// { type : "chat-room/send-message", payload: { roomId: "", text: "" } }

type eventData struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type sendMessagePayload struct {
	RoomID string `json:"roomId"`
	Text   string `json:"text"`
}

func eventsHandler(m *message, s Service) error {
	var ed eventData
	err := json.Unmarshal([]byte(m.Data), &ed)
	if err != nil {
		return err
	}

	if ed.Type == "chat-room/send-message" {
		var smp sendMessagePayload
		if err := json.Unmarshal([]byte(ed.Payload), &smp); err != nil {
			logger.L.Error().Err(err).Msg("data from events is invalid")
			return err
		}

		_, err := strconv.ParseInt(smp.RoomID, 10, 64)
		if err != nil {
			return err
		}
		// _, err = s.SendMessage(m.UserID, smp.Text, roomID)
		if err != nil {
			return err
		}
	}

	return nil
}
