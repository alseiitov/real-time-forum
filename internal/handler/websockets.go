package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/gorilla/websocket"
)

var clients = make(map[int]*client)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2048
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type client struct {
	id        int
	conn      *websocket.Conn
	broadcast chan *WSEvent
}

type WSEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

func (c *client) readPump() {
	defer func() {
		c.conn.Close()
		delete(clients, c.id)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msgText := strings.TrimSpace(strings.Replace(string(messageBytes), "\n", " ", -1))
		if len(msgText) > 0 {
			body := model.Message{
				Message: msgText,
				Date:    time.Now(),
				UserID:  c.id,
				Read:    false,
			}

			event := &WSEvent{Type: "message", Body: body}
			c.broadcast <- event
		}
	}
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case event, ok := <-c.broadcast:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(&event); err != nil {
				log.Println(err)
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.broadcast)
			for i := 0; i < n; i++ {
				if err := c.conn.WriteJSON(<-c.broadcast); err != nil {
					log.Println(err)
					return
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Handler) handleWebSocket(ctx *gorouter.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		conn.WriteJSON(WSEvent{Type: "error", Body: err.Error()})
		return
	}

	c, ok := clients[userID]
	if !ok {
		c = &client{id: userID, conn: conn, broadcast: make(chan *WSEvent, 1)}
		clients[userID] = c
	}
	//TODO: handle multiple connections from one user
	go func() {
		for {
			fmt.Println(clients)
			time.Sleep(2 * time.Second)
		}
	}()
	go c.readPump()
	go c.writePump()
}
