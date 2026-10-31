package link

import (
	"context"
	"errors"
	"testing"
)

type mockStorage struct {
	saveErr       error
	getByCodeErr  error
	getByCodeLink *Link
	savedLinks    []*Link
}

func (m *mockStorage) Save(ctx context.Context, l *Link) error {

	if m.saveErr != nil {
		return m.saveErr
	}
	m.savedLinks = append(m.savedLinks, l)

	return nil
}

func (m *mockStorage) GetByCode(ctx context.Context, code string) (*Link, error) {

	if m.getByCodeErr != nil {
		return nil, m.getByCodeErr
	}
	return m.getByCodeLink, nil
}
func (m *mockStorage) IncrementClicks(ctx context.Context, code string) error {

	return nil
}

var _ Storage = (*mockStorage)(nil)

func TestService_Create_Success(t *testing.T) {

	ctx := context.Background()
	storage := &mockStorage{}
	svc := NewService(storage)

	l, err := svc.Create(ctx, "https://example.com")

	if err != nil {
		t.Fatalf("create: %v:", err)
	}

	if l.Code == "" {
		t.Error("Code is empty")
	}

	if len(l.Code) != codeLength {

		t.Errorf("Code length: got %d, want %d", len(l.Code), codeLength)

	}

	if l.URL != "https://example.com" {
		t.Fatalf("URL: got %q, want %q", l.URL, "https://example.com")
	}

	if len(storage.savedLinks) != 1 {

		t.Errorf("savedLinks: got %q, want 1", len(storage.savedLinks))
	}
}

func TestService_Create_Empty(t *testing.T) {

	ctx := context.Background()
	svc := NewService(&mockStorage{})

	_, err := svc.Create(ctx, "")
	if !errors.Is(err, ErrInvalidURL) {
		t.Errorf("expected ErrInvalidURL, got %v", err)
	}

}

func TestService_Create_InvalidScheme(t *testing.T) {

	ctx := context.Background()
	svc := NewService(&mockStorage{})

	_, err := svc.Create(ctx, "javascript:alert(1)")

	if !errors.Is(err, ErrInvalidURL) {

		t.Errorf("expected ErrInvalidURL, got %v", err)
	}
}

func TestService_Create_RetryOnCollision(t *testing.T) {
	ctx := context.Background()

	callCount := 0
	storage := &mockStorageFunc{
		saveFunc: func(ctx context.Context, l *Link) error {
			callCount++
			if callCount == 1 {
				return ErrCodeExists
			}
			return nil
		},
	}

	svc := NewService(storage)
	l, err := svc.Create(ctx, "https://example.com")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if l == nil {
		t.Fatal("link is nil")
	}
	if callCount != 2 {
		t.Errorf("Save called %d times, want 2", callCount)
	}
}

func TestService_Resole_Success(t *testing.T) {

	ctx := context.Background()

	expected := &Link{Code: "abc123", URL: "https://example.com"}

	storage := &mockStorage{getByCodeLink: expected}
	svc := NewService(storage)

	got, err := svc.Resolve(ctx, "abc123")

	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	if got.Code != expected.Code {

		t.Fatalf("Code: got %q, want %q", got.Code, expected.Code)
	}
}

func TestService_Resole_NotFound(t *testing.T) {

	ctx := context.Background()
	storage := &mockStorage{getByCodeErr: ErrNotFound}
	svc := NewService(storage)

	_, err := svc.Resolve(ctx, "nonexistent")

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

type mockStorageFunc struct {
	saveFunc func(ctx context.Context, l *Link) error
}

func (m *mockStorageFunc) Save(ctx context.Context, l *Link) error {

	if m.saveFunc != nil {
		return m.saveFunc(ctx, l)
	}
	return nil
}

func (m *mockStorageFunc) GetByCode(ctx context.Context, code string) (*Link, error) {

	return nil, ErrNotFound
}

func (m *mockStorageFunc) IncrementClicks(ctx context.Context, code string) error {

	return nil
}
