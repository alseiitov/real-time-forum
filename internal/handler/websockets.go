package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/gorilla/websocket"
)

var clients = make(map[int]*client)

const (
	tokenWait      = 10 * time.Second
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
	mu        sync.Mutex
	broadcast chan *WSEvent
}

type WSEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body,omitempty"`
}

func (h *Handler) clientReadPump(c *client) {
	ticker := time.NewTicker(tokenWait)
	defer ticker.Stop()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	go func() {
		defer func() {
			c.conn.Close()
			delete(clients, c.id)
		}()
		for {
			_, messageBytes, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}

			var event WSEvent
			err = json.Unmarshal(messageBytes, &event)
			if err != nil {
				c.writeJSON(&WSEvent{Type: "error", Body: err.Error()})
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
					c.writeJSON(&WSEvent{Type: "error", Body: err.Error()})
					return
				}
				c.id = sub
				clients[c.id] = c
			case "pongMessage":
				c.conn.SetReadDeadline(time.Now().Add(pongWait))
			default:
				c.writeJSON(&WSEvent{Type: "error", Body: "invalid event type"})
				return
			}
		}
	}()

	<-ticker.C
	if c.id == 0 {
		c.writeJSON(&WSEvent{Type: "error", Body: "no token received"})
		c.conn.Close()
		delete(clients, c.id)
	}
}

func (c *client) writeJSON(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.conn.WriteJSON(data)
	if err != nil {
		log.Println(err)
	}
	return err
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
				c.writeJSON(&WSEvent{Type: "closeMessage"})
				return
			}

			if err := c.writeJSON(&event); err != nil {
				return
			}
			// Add queued chat messages to the current websocket message.
			n := len(c.broadcast)
			for i := 0; i < n; i++ {
				if err := c.writeJSON(<-c.broadcast); err != nil {
					return
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.writeJSON(&WSEvent{Type: "pingMessage"}); err != nil {
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

	// go func() {
	// 	for {
	// 		fmt.Println(clients)
	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()
}
