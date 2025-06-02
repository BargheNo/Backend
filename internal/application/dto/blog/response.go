package blogdto

import "time"

type PostResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      uint      `json:"status"`
	Corporation string    `json:"corporation"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	CoverImage  string    `json:"cover_image"`
	CreatedAt   time.Time `json:"created_at"`
}
