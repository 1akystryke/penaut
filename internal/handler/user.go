package handler

import (
	"io"
	"net/http"
)

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name         string `json:"name"`
		Email        string `json:"email"`
		PasswordHash string `json:"password_hash"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.service.CreateUser(r.Context(), req.Name, req.Email, req.PasswordHash)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, u)
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *Handler) getUserPic(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	user, err := h.service.GetUserbyID(r.Context(), userID)
	picture, err := h.service.GetUserPic(r.Context(), *user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "image/webp")
	defer picture.Close()
	io.Copy(w, picture)
}
