package link

import (
	"context"
	"errors"
	"time"
)

type Link struct {
	Code      string    `json:"code"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	Clicks    int64     `json:"clicks"`
}

var (
	ErrNotFound   = errors.New("link not found")
	ErrCodeExists = errors.New("link code already exists")
	ErrInvalidURL = errors.New("invalid url")
)

type Storage interface {
	Save(ctx context.Context, l *Link) error

	GetByCode(ctx context.Context, code string) (*Link, error)

	IncrementClicks(ctx context.Context, code string) error
}
