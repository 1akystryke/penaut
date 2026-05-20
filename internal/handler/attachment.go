package handler

import (
	"io"
	"net/http"
)

func (h *Handler) GetAttachment(w http.ResponseWriter, r *http.Request) {
	filePath := r.PathValue("file_path")

	picture, err := h.service.GetAttachment(r.Context(), filePath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	defer picture.Close()
	io.Copy(w, picture)
}

func (h *Handler) GetAttachmentsForPost(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("post_id")
	attachments, err := h.service.GetAttachmentsbyPostID(r.Context(), postID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, attachments)

}
