package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/gorilla/websocket"
)

const (
	maxConnsForUser = 3
	tokenWait       = 10 * time.Second
	writeWait       = 10 * time.Second
	pongWait        = 5 * time.Second
	pingPeriod      = (pongWait * 9) / 10
	maxMessageSize  = 256
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var clients = make(map[int]*client)

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

func (h *Handler) runEventsPump() {
	for {
		event, ok := <-h.eventsChan
		if !ok {
			continue
		}
		sendEventToClient(event)
	}
}

func (h *Handler) connReadPump(conn *conn) {
	defer conn.close()

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
			conn.conn.SetReadDeadline(time.Now().Add(pongWait))
		default:
			conn.writeJSON(&model.WSEvent{Type: model.WSEventTypes.Error, Body: "invalid event type"})
			return
		}
	}
}

func (c *conn) ping() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.close()
	}()
	for {
		<-ticker.C
		c.conn.SetWriteDeadline(time.Now().Add(pongWait))
		if err := c.writeJSON(&model.WSEvent{Type: model.WSEventTypes.PingMessage}); err != nil {
			return
		}
	}
}

func (c *conn) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, ok := clients[c.clientID]
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
		delete(clients, client.id)
	}
}

func (c *conn) writeJSON(data interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(data)
}

func sendEventToClient(event *model.WSEvent) {
	client, ok := clients[event.RecipientID]
	if !ok {
		return
	}
	client.send(event)
}

func (c *client) send(event *model.WSEvent) {
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

func (h *Handler) handleWebSocket(ctx *gorouter.Context) {
	ws, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connection := &conn{conn: ws}

	err = h.identifyConn(connection)
	if err != nil {
		connection.writeJSON(&model.WSEvent{Type: model.WSEventTypes.Error, Body: err.Error()})
		connection.conn.Close()
		return
	}

	c, ok := clients[connection.clientID]
	if !ok {
		c = &client{id: connection.clientID}
		clients[connection.clientID] = c
	}

	if len(c.conns) == maxConnsForUser {
		conn := c.conns[0]
		conn.writeJSON(&model.WSEvent{Type: model.WSEventTypes.Error, Body: "too many connections"})
		conn.close()
	}

	go h.connReadPump(connection)
	go connection.ping()

	c.conns = append(c.conns, connection)
}

func (h *Handler) identifyConn(c *conn) error {
	c.conn.SetReadDeadline(time.Now().Add(tokenWait))

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

func UnmarshalEventBody(e *model.WSEvent, v interface{}) error {
	bodyBytes, err := json.Marshal(e.Body.(map[string]interface{}))
	if err != nil {
		return err
	}

	return json.Unmarshal(bodyBytes, &v)
}
