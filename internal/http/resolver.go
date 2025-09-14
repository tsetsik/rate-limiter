package http

import "net/http"

type Handler struct {
}

func NewHttpResolver() *Handler {
	return &Handler{}
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received"))
}
