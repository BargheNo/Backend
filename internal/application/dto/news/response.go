package newsdto

import (
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type AdminNewsResponse struct {
	ID          uint                       `json:"id"`
	Title       string                     `json:"title"`
	Content     string                     `json:"content"`
	Description string                     `json:"description"`
	Status      string                     `json:"status"`
	CoverImage  string                     `json:"coverImage"`
	Author      userdto.CredentialResponse `json:"author"`
}

type PublicNewsResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	CoverImage  string `json:"coverImage"`
}

type NewsStatusesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
