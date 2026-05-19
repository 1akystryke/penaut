package handler

import (
	"net/http"
	"strings"
)

func (h *Handler) createChannel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
		"Bearer ",
	)
	if !ok {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}
	ch, err := h.service.CreateChannel(r.Context(), req.Type, req.Name, token)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, ch)
}

func (h *Handler) CreateDirect(w http.ResponseWriter, r *http.Request) {

	UserID := r.PathValue("user_id")

	token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
		"Bearer ",
	)
	if !ok {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	ch, err := h.service.CreateDirect(r.Context(), UserID, token)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, ch)
}

func (h *Handler) listChannels(w http.ResponseWriter, r *http.Request) {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
		"Bearer ",
	)
	if !ok {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}
	channels, err := h.service.ListChannels(r.Context(), token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, channels)
}
