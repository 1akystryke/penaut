package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"messenger/internal/service"
)

type Handler struct {
	service *service.Service
	hub     *hub
}

func New(service *service.Service) *Handler {
	return &Handler{
		service: service,
		hub:     newHub(),
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /users", h.listUsers)
	mux.HandleFunc("POST /users", h.createUser)
	mux.HandleFunc("GET /channels", h.listChannels)
	mux.HandleFunc("POST /channels", h.createChannel)
	mux.HandleFunc("GET /channels/{channel_id}/members", h.listMembers)
	mux.HandleFunc("POST /channels/{channel_id}/members", h.addMember)
	mux.HandleFunc("GET /channels/{channel_id}/posts", h.listPosts)
	mux.HandleFunc("POST /channels/{channel_id}/posts", h.createPost)
	mux.HandleFunc("POST /auth", h.auth)
	mux.HandleFunc("GET /ws", h.websocket)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func readJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("write json: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
