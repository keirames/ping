package ws

import (
	"chatroom/common/converter"
	"chatroom/logger"
	"encoding/json"
	"fmt"
)

// { type : "chat-room/send-message", payload: { roomId: "", text: "" } }

type eventType struct {
	Type string `json:"type"`
}

type eventPayloadChatRoomSendMessage struct {
	Payload sendMessagePayload `json:"payload"`
}

type sendMessagePayload struct {
	RoomID string `json:"roomId"`
	Text   string `json:"text"`
}

func eventsHandler(m *message, s Service) error {
	var et eventType
	err := json.Unmarshal([]byte(m.Data), &et)
	if err != nil {
		logger.L.Err(err).Msg("cannot unmarshal event type")
		return err
	}

	if et.Type == "chat-room/send-message" {
		var event eventPayloadChatRoomSendMessage
		if err := json.Unmarshal([]byte(m.Data), &event); err != nil {
			logger.L.Err(err).Msg("cannot unmarshal event payload")
			return err
		}

		roomID, err := converter.StringToInt64(event.Payload.RoomID)
		if err != nil {
			logger.L.Err(err).Msg("cannot convert")
			return err
		}

		fmt.Print(roomID)
		// send message
	}

	return nil
}
