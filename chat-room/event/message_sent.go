package event

import (
	"chatroom/ws"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type data struct {
	UserID    int64 `json:"userId"`
	MessageID int64 `json:"messageId"`
	RoomID    int64 `json:"roomId"`
}

func SubscribeMessageSentTopic(hub ws.Hub) {
	Subscribe("message-sent", func(m kafka.Message) {
		var d data
		err := json.Unmarshal(m.Value, &d)
		if err != nil {
			fmt.Println(err)
			return
		}

		hub.SendMessageToClient(d.UserID, d.MessageID)
	})
}
