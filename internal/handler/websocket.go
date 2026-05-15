package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func (h *Handler) websocket(w http.ResponseWriter, r *http.Request) {
	channelID := strings.TrimSpace(r.URL.Query().Get("channel_id"))
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if channelID == "" || token == "" {
		writeError(w, http.StatusBadRequest, errors.New("channel_id and token query params is required"))
		return
	}

	// user, err := h.service.GetUserbyToken(r.Context(), token)
	// if err != nil {
	// 	log.Println("не найден пользователь в вебсокете")
	// 	return
	// }

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
			log.Println(err)
			return
		}

		p, err := h.service.CreatePost(r.Context(), channelID, req.Text, token)
		if err != nil {
			log.Println(err)
			_ = c.writeJSON(map[string]string{"error": err.Error()})
			continue
		}
		log.Println(p.Text)
		h.hub.broadcast(channelID, p)
	}
}
