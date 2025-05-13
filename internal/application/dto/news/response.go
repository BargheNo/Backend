package newsdto

import "github.com/BargheNo/Backend/internal/domain/enum"

type NewsResponse struct {
	ID      uint            `json:"id"`
	Title   string          `json:"title"`
	Content string          `json:"content"`
	Status  enum.NewsStatus `json:"status"`
}
