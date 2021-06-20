package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/alseiitov/gorouter"
	"github.com/alseiitov/real-time-forum/internal/config"
	"github.com/alseiitov/real-time-forum/internal/model"
	"github.com/alseiitov/real-time-forum/internal/service"
	"github.com/alseiitov/real-time-forum/pkg/auth"
	"github.com/gorilla/websocket"
)

type Handler struct {
	clients         map[int]*client
	eventsChan      chan *model.WSEvent
	chatsService    service.Chats
	usersService    service.Users
	tokenManager    auth.TokenManager
	maxConnsForUser int
	maxMessageSize  int64
	tokenWait       time.Duration
	writeWait       time.Duration
	pongWait        time.Duration
	pingPeriod      time.Duration
	upgrader        websocket.Upgrader
}

func NewHandler(eventsChan chan *model.WSEvent, services *service.Services, tokenManager auth.TokenManager, config *config.Conf) *Handler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return &Handler{
		clients:         make(map[int]*client),
		eventsChan:      eventsChan,
		chatsService:    services.Chats,
		usersService:    services.Users,
		tokenManager:    tokenManager,
		upgrader:        upgrader,
		maxConnsForUser: config.Websocket.MaxConnsForUser,
		maxMessageSize:  config.Websocket.MaxMessageSize,
		tokenWait:       config.TokenWait(),
		writeWait:       config.WriteWait(),
		pongWait:        config.PongWait(),
		pingPeriod:      config.PingPeriod(),
	}
}

func (h *Handler) ServeWS(ctx *gorouter.Context) {
	ws, err := h.upgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	connection := &conn{conn: ws}

	err = h.identifyConn(connection)
	if err != nil {
		connection.writeError(err)
		connection.conn.Close()
		return
	}

	c, ok := h.clients[connection.clientID]
	if !ok {
		user, err := h.usersService.GetByID(connection.clientID)
		if err != nil {
			connection.writeError(err)
			connection.conn.Close()
			return
		}

		c = &client{User: user}
		h.clients[connection.clientID] = c
	}

	if len(c.conns) == h.maxConnsForUser {
		conn := c.conns[0]
		conn.writeError(errTooManyConnections)
		h.closeConn(conn)
	}

	go h.connReadPump(connection)
	go h.pingConn(connection)

	c.conns = append(c.conns, connection)
}

func (h *Handler) RunEventsPump() {
	for {
		event, ok := <-h.eventsChan
		if !ok {
			continue
		}
		h.sendEventToClient(event)
	}
}
