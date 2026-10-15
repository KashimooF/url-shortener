package httpapi

import "time"

type CreateLinkRequest struct {
	URL string `json:"url"`
}

type LinkResponse struct {
	Code      string    `json:"code"`
	Short     string    `json:"short"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	Clicks    int64     `json:"clicks"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
