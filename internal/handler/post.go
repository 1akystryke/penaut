package handler

import (
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 MB лимит
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		// если нет файлов — это нормально, продолжаем
		// если превышен размер — возвращаем ошибку
		if err != http.ErrNotMultipart {
			writeError(w, http.StatusBadRequest, err)
			return
		}
	}

	channelID := r.PathValue("channel_id")
	text := r.FormValue("text")

	// теперь безопасно — если файлов нет, вернёт nil slice
	var files []*multipart.FileHeader
	if r.MultipartForm != nil {
		files = r.MultipartForm.File["attachments[]"]
	}

	token, ok := strings.CutPrefix(r.Header.Get("Authorization"),
		"Bearer ",
	)
	if !ok {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	p, err := h.service.CreatePost(r.Context(), channelID, text, token)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			log.Println(err)
			http.Error(w, "Файл не загружен", http.StatusBadRequest)
		}
		defer file.Close()

		// проверь MIME
		buf := make([]byte, 512)
		file.Read(buf)
		mimeType := http.DetectContentType(buf)
		//fileName := fileHeader.Filename
		fileSize := fileHeader.Size
		file.Seek(0, 0) // важно — вернуть указатель в начало перед загрузкой в MinIO

		// генерируй ключ
		h.service.UploadAttachment(r.Context(), file, int(fileSize), mimeType, &p)
		// загружай в MinIO
		// сохраняй в БД: key, fileHeader.Filename, mimeType, fileHeader.Size, postID
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
