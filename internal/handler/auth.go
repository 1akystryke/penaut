package handler

import "net/http"

func (h *Handler) auth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email   string `json:"email"`
		PwdHash string `json:"pwd"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()
	email, pwdHash := req.Email, req.PwdHash
	user, err := h.service.CheckPassword(ctx, email, pwdHash)

	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.service.MakeToken(ctx, user.ID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, token.Token)
}
