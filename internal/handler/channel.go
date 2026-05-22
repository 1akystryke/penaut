package handler

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) createChannel(w http.ResponseWriter, r *http.Request) {

	token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
		"Bearer ",
	)
	if !ok {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	r.ParseMultipartForm(32 << 20)

	// 2. Получаем текстовые поля
	Type := r.FormValue("type")
	Name := r.FormValue("name")

	// 3. Получаем файл
	file, header, err := r.FormFile("pic")

	if err != nil {
		log.Println(err)
		http.Error(w, "Файл не загружен", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ch, err := h.service.CreateChannel(r.Context(), Type, Name, token)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	_, err = h.service.UploadChannelPic(r.Context(), &ch, file, int(header.Size), header.Header.Get("Content-Type"))

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

func (h *Handler) getChannelPic(w http.ResponseWriter, r *http.Request) {
	channelID := r.PathValue("channel_id")
	channel, err := h.service.GetChannelbyID(r.Context(), channelID)
	picture, err := h.service.GetChannelPic(r.Context(), *channel)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	defer picture.Close()
	io.Copy(w, picture)
}
