package handler

import (
	"sync"

	"github.com/gorilla/websocket"
)

type hub struct {
	mu       sync.RWMutex
	channels map[string]map[*client]bool
}

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (c *client) writeJSON(message any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(message)
}

func newHub() *hub {
	return &hub{channels: make(map[string]map[*client]bool)}
}

func (h *hub) add(channelID string, c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.channels[channelID] == nil {
		h.channels[channelID] = make(map[*client]bool)
	}
	h.channels[channelID][c] = true
}

func (h *hub) remove(channelID string, c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.channels[channelID], c)
	if len(h.channels[channelID]) == 0 {
		delete(h.channels, channelID)
	}
}

func (h *hub) broadcast(channelID string, message any) {
	h.mu.RLock()
	clients := make([]*client, 0, len(h.channels[channelID]))
	for c := range h.channels[channelID] {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		if err := c.writeJSON(message); err != nil {
			h.remove(channelID, c)
			_ = c.conn.Close()
		}
	}
}
