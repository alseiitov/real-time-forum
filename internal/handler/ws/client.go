package ws

import (
	"sync"
	"time"

	"github.com/alseiitov/real-time-forum/internal/model"
)

type client struct {
	id    int
	conns []*conn
	mu    sync.Mutex
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

		if err := conn.writeJSON(&event); err != nil {
			h.closeConn(conn)
			continue
		}
	}
}
