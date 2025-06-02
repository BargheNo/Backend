package blogdto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
)

type CorporationPostResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      uint      `json:"status"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	CoverImage  string    `json:"cover_image"`
	CreatedAt   time.Time `json:"created_at"`
}

type GeneralPostResponse struct {
	ID          uint                                         `json:"id"`
	Title       string                                       `json:"title"`
	Description string                                       `json:"description"`
	Status      uint                                         `json:"status"`
	Content     string                                       `json:"content"`
	Corporation corporationdto.CorporationCredentialResponse `json:"corporation"`
	CoverImage  string                                       `json:"cover_image"`
	CreatedAt   time.Time                                    `json:"created_at"`
}
