package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"proektUrl/internal/link"
)

type Handler struct {
	svc    *link.Service
	host   string
	logger *slog.Logger
}

func NewHandler(svc *link.Service, host string, logger *slog.Logger) *Handler {

	return &Handler{
		svc:    svc,
		host:   host,
		logger: logger,
	}
}

func (h *Handler) CreateLink(w http.ResponseWriter, r *http.Request) {

	var req CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	l, err := h.svc.Create(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, link.ErrInvalidURL) {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.logger.Error("create link failed", "err", err, "url", req.URL)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := LinkResponse{
		Code:      l.Code,
		Short:     h.host + "/" + l.Code,
		URL:       l.URL,
		CreatedAt: l.CreatedAt,
		Clicks:    l.Clicks,
	}

	h.writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {

	code := r.PathValue("code")

	l, err := h.svc.Resolve(r.Context(), code)

	if err != nil {

		if errors.Is(err, link.ErrNotFound) {
			h.writeError(w, http.StatusNotFound, "link not found")
			return
		}
		h.logger.Error("resolve failed", "err", err, "code", code)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	http.Redirect(w, r, l.URL, http.StatusFound)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, v any) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {

		h.logger.Error("encode response", "err", err)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, msg string) {

	h.writeJSON(w, status, ErrorResponse{Error: msg})
}
