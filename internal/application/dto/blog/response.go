package blogdto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type CorporationPostResponse struct {
	ID          uint                       `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Status      string                     `json:"status"`
	Content     string                     `json:"content"`
	Author      userdto.CredentialResponse `json:"author"`
	CoverImage  string                     `json:"coverImage"`
	CreatedAt   time.Time                  `json:"createdAt"`
	LikeCount   int                        `json:"likeCount"`
}

type GeneralPostResponse struct {
	ID          uint                                         `json:"id"`
	Title       string                                       `json:"title"`
	Description string                                       `json:"description"`
	Content     string                                       `json:"content"`
	Corporation corporationdto.CorporationCredentialResponse `json:"corporation"`
	CoverImage  string                                       `json:"coverImage"`
	CreatedAt   time.Time                                    `json:"createdAt"`
	LikeCount   int                                          `json:"likeCount"`
}

type GetBlogEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
