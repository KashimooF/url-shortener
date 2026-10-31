package link

import (
	"context"
	"errors"
	"testing"
)

func TestMemoryStorage_SaveAndGet(t *testing.T) {

	ctx := context.Background()
	storage := NewMemoryStorage()

	l := &Link{
		Code: "abc123",
		URL:  "https://example.com",
	}
	//Save

	if err := storage.Save(ctx, l); err != nil {

		t.Fatalf("Save: %v", err)
	}

	//GetByCode

	got, err := storage.GetByCode(ctx, "abc123")

	if err != nil {
		t.Fatalf("GetByCode: %v", err)
	}

	if got.Code != l.Code {

		t.Fatalf("Code: got %q, want %q", got.Code, l.Code)
	}

	if got.URL != l.URL {

		t.Fatalf("URL: got %q, want %q", got.URL, l.URL)
	}
}

func TestMemoryStorage_Save_Duplicate(t *testing.T) {

	ctx := context.Background()

	storage := NewMemoryStorage()

	l := &Link{
		Code: "abc123",
		URL:  "https://example.com",
	}

	if err := storage.Save(ctx, l); err != nil {

		t.Fatalf("first save: %v", err)
	}

	err := storage.Save(ctx, l)

	if !errors.Is(err, ErrCodeExists) {
		t.Fatalf("expected ErrCodeExist, got %v", err)
	}
}

func TestMemoryStorage_GetByCode_NotFound(t *testing.T) {

	ctx := context.Background()
	storage := NewMemoryStorage()

	_, err := storage.GetByCode(ctx, "nonexistent")

	if !errors.Is(err, ErrNotFound) {

		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStorage_IncrementClicks(t *testing.T) {

	ctx := context.Background()

	storage := NewMemoryStorage()

	l := &Link{
		Code:   "abc123",
		URL:    "https://example.com",
		Clicks: 0,
	}

	if err := storage.Save(ctx, l); err != nil {

		t.Fatalf("setup save: %v", err)
	}

	//два раза увеличиваем

	if err := storage.IncrementClicks(ctx, "abc123"); err != nil {

		t.Fatalf("IncrementClicks 1: %v", err)
	}

	if err := storage.IncrementClicks(ctx, "abc123"); err != nil {

		t.Fatalf("IncrementClicks 2: %v", err)
	}

	got, err := storage.GetByCode(ctx, "abc123")

	if err != nil {
		t.Fatalf("GetByCode: %v", err)
	}

	if got.Clicks != 2 {
		t.Fatalf("Clicks: got %d, want 2", got.Clicks)
	}
}
