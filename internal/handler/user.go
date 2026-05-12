package handler

import "net/http"

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.service.CreateUser(r.Context(), req.Name, req.Email)
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
