package files

import (
	"net/http"
	"quizzly/pkg/files"
	"quizzly/pkg/logger"
)

const (
	pathValueFilename = "filename"
)

type GetFileHandler struct {
	file files.Manager
	log  logger.Logger
}

func NewGetFileHandler(file files.Manager, log logger.Logger) *GetFileHandler {
	return &GetFileHandler{
		file: file,
		log:  log,
	}
}

func (h *GetFileHandler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue(pathValueFilename)
		if filename == "" {
			return
		}

		buffer, err := h.file.Get(r.Context(), filename)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = w.Write(buffer.Bytes())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			h.log.Error("handle request error", err)
		}
	}
}
