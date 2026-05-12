package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func (h *Handler) websocket(w http.ResponseWriter, r *http.Request) {
	channelID := strings.TrimSpace(r.URL.Query().Get("channel_id"))
	if channelID == "" {
		writeError(w, http.StatusBadRequest, errors.New("channel_id query param is required"))
		return
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	c := &client{conn: conn}
	h.hub.add(channelID, c)
	defer h.hub.remove(channelID, c)

	for {
		var req struct {
			Text string `json:"text"`
		}
		if err := conn.ReadJSON(&req); err != nil {
			return
		}

		p, err := h.service.CreatePost(r.Context(), channelID, req.Text)
		if err != nil {
			_ = c.writeJSON(map[string]string{"error": err.Error()})
			continue
		}
		h.hub.broadcast(channelID, p)
	}
}
