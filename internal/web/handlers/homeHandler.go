package handlers

import "net/http"

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to Order Service!"))
}

func NewHomeHandler() *homeHandler {
	return &homeHandler{}
}
