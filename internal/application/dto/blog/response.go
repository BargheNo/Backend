package blogdto

import "time"

type PostResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Corporation string    `json:"corporation"`
	Author      string    `json:"author"`
	CoverImage  string    `json:"cover_image"`
	CreatedAt   time.Time `json:"created_at"`
}

type PostDetailsResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Corporation string    `json:"corporation"`
	Content     string    `json:"content"`
	CoverImage  string    `json:"cover_image"`
	CreatedAt   time.Time `json:"created_at"`
}
