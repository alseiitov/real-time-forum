package handler

import (
	"encoding/json"
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
	tokenWait      = 60 * time.Second
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

func (h *Handler) clientReadPump(c *client) {
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

		var event WSEvent
		err = json.Unmarshal(messageBytes, &event)
		if err != nil {
			if err := c.conn.WriteJSON(&WSEvent{Type: "error", Body: err.Error()}); err != nil {
				log.Println(err)
				return
			}
			return
		}

		switch event.Type {
		case "message":
			if c.id != 0 {
				msgText := strings.TrimSpace(strings.Replace(fmt.Sprintf("%s", event.Body), "\n", " ", -1))
				if len(msgText) > 0 {
					body := model.Message{
						Message: msgText,
						Date:    time.Now(),
						UserID:  c.id,
						Read:    false,
					}
					c.broadcast <- &WSEvent{Type: "message", Body: body}
				}
			}
		case "token":
			token := fmt.Sprintf("%s", event.Body)
			sub, _, err := h.tokenManager.Parse(token)
			if err != nil {
				if err := c.conn.WriteJSON(&WSEvent{Type: "error", Body: err.Error()}); err != nil {
					log.Println(err)
					return
				}
				return
			}
			c.id = sub
			clients[c.id] = c
		default:
			if err := c.conn.WriteJSON(&WSEvent{Type: "error", Body: "invalid event type"}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (h *Handler) clientWritePump(c *client) {
	ticker := time.NewTicker(writeWait)
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

	//TODO: handle multiple connections from one user
	c := &client{conn: conn, broadcast: make(chan *WSEvent, 16)}

	go h.clientReadPump(c)
	go h.clientWritePump(c)

	go func() {
		for {
			fmt.Println(clients)
			time.Sleep(5 * time.Second)
		}
	}()
}
