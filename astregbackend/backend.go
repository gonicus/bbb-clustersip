package main

import (
	"net/http"
)

type Backend interface {
	Single(w http.ResponseWriter, r *http.Request)
	Multi(w http.ResponseWriter, r *http.Request)
}

type RealtimeHandler struct {
	mux     *http.ServeMux
	backend Backend
}

func (h *RealtimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, pattern := h.mux.Handler(r)
	if pattern != "" && r.Method != "POST" {
		http.Error(w, "405 Method Not Allowed", 405)
		return
	}

	err := r.ParseForm()
	r.Body.Close()
	if err != nil {
		http.Error(w, "", 400)
	}

	w.Header().Set("Content-Type", "text/plain")
	handler.ServeHTTP(w, r)
}

func NewRealtimeHandler(pattern string, backend Backend) (h *RealtimeHandler) {
	h = &RealtimeHandler{}
	h.backend = backend
	h.mux = http.NewServeMux()

	h.mux.HandleFunc(pattern+"single", h.backend.Single)
	h.mux.HandleFunc(pattern+"multi", h.backend.Multi)

	h.mux.HandleFunc(pattern+"require", func(w http.ResponseWriter, r *http.Request) {
		// Ignore request
		w.Write([]byte("0"))
	})

	http.Handle(pattern, h)
	return
}
