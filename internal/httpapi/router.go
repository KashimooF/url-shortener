package httpapi

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", h.CreateLink)
	mux.HandleFunc("GET /{code}", h.Redirect)
	return mux
}
