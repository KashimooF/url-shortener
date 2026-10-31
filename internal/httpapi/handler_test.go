package httpapi

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"proektUrl/internal/link"
	"testing"
)

func newTestHandler() *Handler {
	storage := link.NewMemoryStorage()

	svc := link.NewService(storage)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	return NewHandler(svc, "http://test", logger)
}

func TestCreateLink_Success(t *testing.T) {

	h := newTestHandler()

	body := []byte(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	if w.Code != http.StatusCreated {

		t.Fatalf("status: got %d, want %d", w.Code, http.StatusCreated)

	}

	var resp LinkResponse

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {

		t.Fatalf("decode response: %v", err)
	}

	if resp.Code == "" {

		t.Error("Code is empty")
	}

	if resp.URL != "https://example.com" {
		t.Errorf("URL: got %q, wanr %q", resp.URL, "https://example.com")
	}

	if resp.Short != "http://test/"+resp.Code {
		t.Errorf("short: got %q, want %q", resp.Short, "http://test/"+resp.Code)
	}
}

func TestCreateLink_InvalidJSON(t *testing.T) {

	h := newTestHandler()

	body := []byte(`{not json}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)

	}

}

func TestCreateLink_EmptyURL(t *testing.T) {

	h := newTestHandler()

	body := []byte(`{"url":""}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateLink_InvalidScheme(t *testing.T) {

	h := newTestHandler()

	body := []byte(`{"url":"javascript:alert(1)"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreateLink(w, req)

	if w.Code != http.StatusBadRequest {

		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)

	}

}

func TestRedirect_Success(t *testing.T) {
	h := newTestHandler()

	// Сначала создаём ссылку через POST
	body := []byte(`{"url":"https://example.com"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	createW := httptest.NewRecorder()
	h.CreateLink(createW, createReq)

	var created LinkResponse
	if err := json.NewDecoder(createW.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	// Теперь GET /{code}
	req := httptest.NewRequest(http.MethodGet, "/"+created.Code, nil)
	req.SetPathValue("code", created.Code)
	w := httptest.NewRecorder()

	h.Redirect(w, req)

	if w.Code != http.StatusFound {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusFound)
	}
	if loc := w.Header().Get("Location"); loc != "https://example.com" {
		t.Errorf("Location: got %q, want %q", loc, "https://example.com")
	}
}

func TestRedirect_NotFound(t *testing.T) {

	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	req.SetPathValue("code", "nonexistent")
	w := httptest.NewRecorder()

	h.Redirect(w, req)

	if w.Code != http.StatusNotFound {

		t.Fatalf("status: got %d, want %d", w.Code, http.StatusNotFound)
	}
}
