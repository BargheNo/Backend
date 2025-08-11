package newsdto

import (
	"time"

	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type AdminNewsResponse struct {
	ID          uint                       `json:"id"`
	CreatedAt   time.Time                  `json:"createdAt"`
	Title       string                     `json:"title"`
	Content     string                     `json:"content"`
	Description string                     `json:"description"`
	Status      string                     `json:"status"`
	CoverImage  string                     `json:"coverImage"`
	Author      userdto.CredentialResponse `json:"author"`
	TotalLike   int                        `json:"totalLikes"`
}

type PublicNewsResponse struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	CoverImage  string    `json:"coverImage"`
	TotalLike   int       `json:"totalLikes"`
}

type NewsEnumResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
