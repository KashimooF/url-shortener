package link

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

const (
	codeLength  = 6
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxAttempts = 5
)

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {

	return &Service{storage: storage}
}

func (s *Service) Create(ctx context.Context, url string) (*Link, error) {

	if url == "" {
		return nil, fmt.Errorf("%w: empty", ErrInvalidURL)

	}

	for i := 0; i < maxAttempts; i++ {

		code, err := generateCode(codeLength)

		if err != nil {
			return nil, fmt.Errorf("generate code: %w", err)
		}

		l := &Link{
			Code:      code,
			URL:       url,
			CreatedAt: time.Now().UTC(),
			Clicks:    0,
		}
		err = s.storage.Save(ctx, l)

		if err == nil {

			return l, nil
		}
		if errors.Is(err, ErrCodeExists) {
			continue
		}
		return nil, fmt.Errorf("save link: %w", err)
	}
	return nil, errors.New("failed to generate unique code")
}

func (s *Service) Resolve(ctx context.Context, code string) (*Link, error) {

	l, err := s.storage.GetByCode(ctx, code)

	if err != nil {

		return nil, err
	}

	if err := s.storage.IncrementClicks(ctx, code); err != nil {

		_ = err
	}

	return l, nil
}

func generateCode(n int) (string, error) {
	b := make([]byte, n)
	max := big.NewInt(int64(len(alphabet)))

	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[idx.Int64()]
	}
	return string(b), nil
}
