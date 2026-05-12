package handler

import "net/http"

func (h *Handler) createChannel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type string `json:"type"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ch, err := h.service.CreateChannel(r.Context(), req.Type)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, ch)
}

func (h *Handler) listChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := h.service.ListChannels(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, channels)
}
