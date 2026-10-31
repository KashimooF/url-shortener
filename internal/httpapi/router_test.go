package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_Routes(t *testing.T) {

	h := newTestHandler()
	router := NewRouter(h)

	t.Run("POST /shorten", func(t *testing.T) {
		body := strings.NewReader(`{"url":"https://example.com"}`)
		req := httptest.NewRequest(http.MethodPost, "/shorten", body)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusCreated)

		}
	})

	t.Run("GET /{code} not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {

			t.Errorf("status: got %d, wantr %d", w.Code, http.StatusNotFound)

		}
	})

	t.Run("GET /shorten is not POST /shorten", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusNotFound)
		}

	})

	t.Run("DELETE not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/shorten", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("status: got %d, want %d", w.Code, http.StatusMethodNotAllowed)

		}
	})

}
