package ws

import (
	"sync"
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type client struct {
	conns []*conn
	mu    sync.Mutex
	model.User
}

func (h *Handler) sendEventToClient(event *model.WSEvent) {
	client, ok := h.clients[event.RecipientID]
	if !ok {
		return
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	for i := 0; i < len(client.conns); i++ {
		conn := client.conns[i]
		conn.conn.SetWriteDeadline(time.Now().Add(h.writeWait))

		err := conn.writeJSON(&event)
		if err != nil {
			h.closeConn(conn)
			continue
		}
	}
}
