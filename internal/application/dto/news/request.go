package newsdto

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreateNewsRequest struct {
	Title       string
	Content     string
	Description string
	AuthorID    uint
	Status      enum.NewsStatus
	CoverImage  *multipart.FileHeader
}

type EditNewsRequest struct {
	NewsID      uint
	AuthorID    uint
	Title       *string
	Content     *string
	Description *string
	Status      uint
	CoverImage  *multipart.FileHeader
}

type EditNewsStatusRequest struct {
	NewsID   uint
	AuthorID uint
	Status   uint
}

type DeleteNewsRequest struct {
	NewsIDs  []uint
	AuthorID uint
}

type GetAdminNewsListRequest struct {
	Status uint
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type GetPublicNewsListRequest struct {
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type AddNewsMediaRequest struct {
	NewsID   uint
	AuthorID uint
	Media    *multipart.FileHeader
}

type AccessMediaRequest struct {
	NewsID   uint
	AuthorID uint
	MediaID  uint
	UserType enum.UserType
}

type GetNewsByCustomer struct {
	NewsID uint
	UserID uint
}

type SearchNewsRequest struct {
	Query  string
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}
