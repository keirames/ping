package ws

import (
	"bytes"
	"chatroom/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	readWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type client struct {
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
	defer func() {
		c.conn.Close()
	}()

	msg, ok := <-c.send

	err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
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

func (c *client) readBump() {
	defer func() {
		c.hub.unsubscribe <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readWait))

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				fmt.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.events <- message
	}
}

func Serve(h *hub, w http.ResponseWriter, r *http.Request) {
	// TODO: DEV only
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.L.Error().Err(err).Msg("Cannot upgrade http")
		return
	}

	c := &client{h, conn, make(chan []byte)}
	c.hub.subscribe <- c

	go c.writeBump()
	go c.readBump()
}
