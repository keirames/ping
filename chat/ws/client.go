package ws

import (
	"bytes"
	"context"
	"fmt"
	"main/internal/auth"
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

			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				logger.L.Err(err).Msg("SetWriteDeadline fail")
				return
			}

			err = c.conn.WriteMessage(websocket.PingMessage, nil)
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
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					logger.L.Err(err).Msg("fail to close message")
				}

				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.L.Err(err).Msg("fail to call next writer")
				return
			}
			_, err = w.Write(msg)
			if err != nil {
				logger.L.Err(err).Msg("fail to write")
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
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		logger.L.Err(err).Msg("SetReadDeadline fail")
		return
	}

	c.conn.SetPongHandler(func(string) error {
		logger.L.Info().Msg("pong from user received -> reset read deadline")

		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			logger.L.Err(err).Msg("SetReadDeadline fail")
			return err
		}

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
				logger.L.Error().Msg("UnexpectedCloseError")
			}
			break
		}

		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		c.hub.events <- &message{UserID: id, Data: string(msg)}
	}
}

func Serve(
	ctx context.Context, h *hub, w http.ResponseWriter, r *http.Request,
) error {
	uc, err := auth.GetUser(r.Context())
	if err != nil {
		return err
	}

	// if config.C.Env == "dev" {
	// 	upgrader.CheckOrigin = func(r *http.Request) bool {
	// 		return true
	// 	}
	// }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L.Error().Err(err).Msg("Cannot upgrade http")
		return fmt.Errorf("Cannot upgrade http")
	}

	c := &client{
		id:   uc.ID,
		hub:  h,
		conn: conn,
		send: make(chan []byte),
	}
	c.hub.subscribe <- c

	go c.writeBump()
	go c.readBump(uc.ID)

	return nil
}
