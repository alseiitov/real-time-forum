package service

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

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

type WSEvent struct {
	Type string      `json:"type"`
	Body interface{} `json:"body,omitempty"`
}

type WebsocketsService struct {
	maxConnsForUser int
	maxMessageSize  int
	tokenWait       time.Duration
	writeWait       time.Duration
	pongWait        time.Duration
	pingPeriod      time.Duration
}

func NewWebsocketsService(maxConnsForUser, maxMessageSize int, tokenWait, writeWait, pongWait, pingPeriod time.Duration) *WebsocketsService {
	return &WebsocketsService{
		maxConnsForUser: maxConnsForUser,
		maxMessageSize:  maxMessageSize,
		tokenWait:       tokenWait,
		writeWait:       writeWait,
		pongWait:        pongWait,
		pingPeriod:      pingPeriod,
	}
}

func (s *WebsocketsService) AddClient(c *client) {

}

func (s *WebsocketsService) SendToClient(clientID int, event *WSEvent) {

}
