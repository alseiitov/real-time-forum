package handler

import (
	"encoding/json"
	"errors"
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
	maxConnsForUser = 3
	tokenWait       = 10 * time.Second
	writeWait       = 10 * time.Second
	pongWait        = 60 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	maxMessageSize  = 2048
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type client struct {
	id    int
	conns []*conn
	mu    sync.Mutex
}

type conn struct {
	clientID int
	conn     *websocket.Conn
	mu       sync.Mutex
}
type WSEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body,omitempty"`
}

func (h *Handler) connReadPump(conn *conn) {
	client := clients[conn.clientID]
	defer func() {
		conn.close()
	}()

	conn.conn.SetReadLimit(maxMessageSize)
	conn.conn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		_, messageBytes, err := conn.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		var event WSEvent
		err = json.Unmarshal(messageBytes, &event)
		if err != nil {
			conn.writeJSON(&WSEvent{Type: "error", Body: err.Error()})
			return
		}

		switch event.Type {
		case "message":
			if conn.clientID != 0 {
				msgText := strings.TrimSpace(strings.Replace(fmt.Sprintf("%s", event.Body), "\n", " ", -1))
				if len(msgText) > 0 {
					body := model.Message{
						Message: msgText,
						Date:    time.Now(),
						UserID:  conn.clientID,
						Read:    false,
					}
					client.send(&WSEvent{Type: "message", Body: body})
				}
			}
		case "pongMessage":
			conn.conn.SetReadDeadline(time.Now().Add(pongWait))
		default:
			conn.writeJSON(&WSEvent{Type: "error", Body: "invalid event type"})
			return
		}
	}
}

func (h *Handler) pingConn(conn *conn) {
	ticker := time.NewTicker(writeWait)
	defer func() {
		ticker.Stop()
		conn.close()
	}()
	for {
		<-ticker.C
		conn.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := conn.writeJSON(&WSEvent{Type: "pingMessage"}); err != nil {
			return
		}
	}
}

func (c *client) send(event *WSEvent) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < len(c.conns); i++ {
		conn := c.conns[i]
		conn.conn.SetWriteDeadline(time.Now().Add(writeWait))

		if err := conn.writeJSON(&event); err != nil {
			conn.close()
			continue
		}
	}
}

func (c *conn) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	client := clients[c.clientID]

	for i := 0; i < len(client.conns); i++ {
		if client.conns[i] == c {
			client.conns[i].conn.Close()
			client.conns = append(client.conns[:i], client.conns[i+1:]...)
			break
		}
	}

	if len(client.conns) == 0 {
		delete(clients, client.id)
	}
}

func (h *Handler) identifyConn(c *conn) (int, error) {
	c.conn.SetReadDeadline(time.Now().Add(tokenWait))

	var event WSEvent
	_, messageBytes, err := c.conn.ReadMessage()
	if err != nil {
		return -1, errors.New("no token received")
	}

	err = json.Unmarshal(messageBytes, &event)
	if err != nil {
		return -1, err
	}

	token := fmt.Sprintf("%s", event.Body)
	sub, _, err := h.tokenManager.Parse(token)
	if err != nil {
		return -1, err
	}

	return sub, nil
}

func (c *conn) writeJSON(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(data)
}

func (h *Handler) handleWebSocket(ctx *gorouter.Context) {
	c, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connection := &conn{conn: c}

	id, err := h.identifyConn(connection)
	if err != nil {
		connection.writeJSON(&WSEvent{Type: "error", Body: err.Error()})
		connection.conn.Close()
		return
	}

	connection.clientID = id

	if id != 0 {
		c, ok := clients[id]
		if !ok {
			c = &client{id: id}
			clients[id] = c
		}

		if len(clients[id].conns) == maxConnsForUser {
			conn := clients[id].conns[0]
			conn.writeJSON(&WSEvent{Type: "error", Body: "too many connections"})
			conn.close()
		}

		go h.connReadPump(connection)
		go h.pingConn(connection)

		clients[id].conns = append(clients[id].conns, connection)
	}
}

func logConns() {
	for {
		fmt.Println(len(clients), "clients connected")
		for _, client := range clients {
			fmt.Printf("client %v have %v connections\n", client.id, len(client.conns))
		}
		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}
