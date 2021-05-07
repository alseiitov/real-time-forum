package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/gorilla/websocket"
)

var chats = map[string]*Chat{}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2048
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Chat struct {
	clients    map[*Client]bool // Registered clients.
	broadcast  chan WSEvent     // Inbound messages from the clients.
	register   chan *Client     // Register requests from the clients.
	unregister chan *Client     // Unregister requests from clients.
}

type Client struct {
	id   int
	chat *Chat
	conn *websocket.Conn // The websocket connection.
	send chan WSEvent    // Buffered channel of outbound messages.
}

type WSEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}

func newChat() *Chat {
	return &Chat{
		broadcast:  make(chan WSEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (c *Chat) run() {
	for {
		select {
		case client := <-c.register:
			c.clients[client] = true
		case client := <-c.unregister:
			if _, ok := c.clients[client]; ok {
				delete(c.clients, client)
				close(client.send)
			}
		case message := <-c.broadcast:
			for client := range c.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(c.clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.chat.unregister <- c
		c.conn.Close()
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

			event := WSEvent{Type: "message", Body: body}
			c.chat.broadcast <- event
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(&event); err != nil {
				log.Println(err)
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if err := c.conn.WriteJSON(<-c.send); err != nil {
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

func (h *Handler) handleChatWebSocket(ctx *gorouter.Context) {
	conn, err := upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	chatID, err := ctx.GetStringParam("chat_id")
	if err != nil {
		conn.WriteJSON(WSEvent{Type: "error", Body: err.Error()})
		return
	}

	userID, err := ctx.GetIntParam("sub")
	if err != nil {
		conn.WriteJSON(WSEvent{Type: "error", Body: err.Error()})
		return
	}

	chat, ok := chats[chatID]
	if !ok {
		chat = newChat()
		chats[chatID] = chat
	}
	go chat.run()

	client := &Client{id: userID, chat: chat, conn: conn, send: make(chan WSEvent, 256)}
	client.chat.register <- client
	go client.writePump()
	go client.readPump()
}
