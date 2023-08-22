package ws

import (
	"bytes"
	"fmt"
	"main/config"
	"main/logger"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time allowed to write a message to the peer.
	writeWait = 2 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type client struct {
	// Specific client ID, using user's id
	id int64

	hub *hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *client) writeBump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			logger.L.Info().Msg("ticker end -> ping connected user")

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				logger.L.Err(err).Msg("write message fail")
				return
			}

		case msg, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.L.Err(err).Msg("cannot write deadline")
				return
			}

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, err = w.Write(msg)
			if err != nil {
				return
			}
		}
	}
}

func (c *client) readBump(id int64) {
	defer func() {
		c.hub.unsubscribe <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		logger.L.Info().Msg("pong from user received -> reset read deadline")

		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			logger.L.Err(err).Msg("cannot read message")
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		c.hub.events <- &message{UserID: id, Data: string(msg)}
	}
}

func Serve(h *hub, w http.ResponseWriter, r *http.Request) {
	// userID, _ := middlewares.GetUserID(r.Context())
	// TODO: some aid
	var userID int64

	if config.C.Env == "dev" {
		upgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L.Error().Err(err).Msg("Cannot upgrade http")
		return
	}

	c := &client{userID, h, conn, make(chan []byte)}
	c.hub.subscribe <- c

	go c.writeBump()
	go c.readBump(userID)
}
