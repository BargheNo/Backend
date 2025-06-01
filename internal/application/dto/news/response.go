package newsdto

import "github.com/BargheNo/Backend/internal/domain/enum"

type NewsResponse struct {
	ID         uint            `json:"id"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	Status     enum.NewsStatus `json:"status"`
	CoverImage string          `json:"cover_image"`
}

type NewsStatusesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
