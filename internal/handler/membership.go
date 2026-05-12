package handler

import "net/http"

func (h *Handler) addMember(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("channel_id")
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	m, err := h.service.AddMember(r.Context(), req.UserID, channelID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (h *Handler) listMembers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListMembers(r.Context(), r.PathValue("channel_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}
