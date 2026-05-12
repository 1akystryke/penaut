package handler

import "net/http"

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("channel_id")
	var req struct {
		Text string `json:"text"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.service.CreatePost(r.Context(), channelID, req.Text)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	h.hub.broadcast(channelID, p)
	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) listPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.ListPosts(r.Context(), r.PathValue("channel_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, posts)
}
