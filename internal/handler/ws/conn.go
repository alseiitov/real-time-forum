package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/gorilla/websocket"
)

type conn struct {
	clientID int
	conn     *websocket.Conn
	mu       sync.Mutex
}

func (h *Handler) identifyConn(c *conn) error {
	c.conn.SetReadDeadline(time.Now().Add(h.tokenWait))

	var event model.WSEvent
	_, messageBytes, err := c.conn.ReadMessage()
	if err != nil {
		return errors.New("no token received")
	}

	err = json.Unmarshal(messageBytes, &event)
	if err != nil {
		return err
	}

	token := fmt.Sprintf("%s", event.Body)
	sub, _, err := h.tokenManager.Parse(token)
	if err != nil {
		return err
	}

	c.clientID = sub
	return nil
}

func (h *Handler) connReadPump(conn *conn) {
	defer h.closeConn(conn)

	conn.conn.SetReadLimit(h.maxMessageSize)
	conn.conn.SetReadDeadline(time.Now().Add(h.pongWait))

	for {
		_, messageBytes, err := conn.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		var event model.WSEvent
		err = json.Unmarshal(messageBytes, &event)
		if err != nil {
			conn.writeJSON(&model.WSEvent{Type: model.WSEventTypes.Error, Body: err.Error()})
			return
		}

		switch event.Type {
		case model.WSEventTypes.Message:
			err := h.messageHandler(conn.clientID, &event)
			if err != nil {
				conn.writeJSON(&model.WSEvent{Type: model.WSEventTypes.Error, Body: err.Error()})
				return
			}
		case model.WSEventTypes.PongMessage:
			conn.conn.SetReadDeadline(time.Now().Add(h.pongWait))
		default:
			conn.writeJSON(&model.WSEvent{Type: model.WSEventTypes.Error, Body: "invalid event type"})
			return
		}
	}
}

func (h *Handler) pingConn(c *conn) {
	ticker := time.NewTicker(h.pingPeriod)
	defer func() {
		ticker.Stop()
		h.closeConn(c)
	}()
	for {
		<-ticker.C
		c.conn.SetWriteDeadline(time.Now().Add(h.pongWait))
		if err := c.writeJSON(&model.WSEvent{Type: model.WSEventTypes.PingMessage}); err != nil {
			return
		}
	}
}

func (h *Handler) closeConn(c *conn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, ok := h.clients[c.clientID]
	if !ok {
		return
	}

	for i := 0; i < len(client.conns); i++ {
		if client.conns[i] == c {
			client.conns[i].conn.Close()
			client.conns = append(client.conns[:i], client.conns[i+1:]...)
			break
		}
	}

	if len(client.conns) == 0 {
		delete(h.clients, client.id)
	}
}

func (c *conn) writeJSON(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(data)
}
