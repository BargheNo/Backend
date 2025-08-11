package blogdto

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreatePostRequest struct {
	Title         string
	Content       string
	Description   string
	AuthorID      uint
	CorporationID uint
	CoverImage    *multipart.FileHeader
	Status        enum.PostStatus
}

type EditPostRequest struct {
	PostID        uint
	AuthorID      uint
	CorporationID uint
	Title         *string
	Content       *string
	Description   *string
	CoverImage    *multipart.FileHeader
	Status        *uint
}

type GetPublicPostsRequest struct {
	Offset int
	Limit  int
	SortBy uint
	Asc    bool
}

type GetPublicCorporationPostsRequest struct {
	CorporationID uint
	Offset        int
	Limit         int
	SortBy        uint
	Asc           bool
}

type GetCorporationPostsRequest struct {
	UserID        uint
	CorporationID uint
	Status        uint
	Offset        int
	Limit         int
	SortBy        uint
	Asc           bool
}

type DeletePostRequest struct {
	PostIDs       []uint
	AuthorID      uint
	CorporationID uint
}

type AddPostMediaRequest struct {
	PostID        uint
	AuthorID      uint
	Media         *multipart.FileHeader
	CorporationID uint
}

type AccessPostMediaRequest struct {
	PostID        uint
	UserID        uint
	MediaID       uint
	CorporationID uint
	UserType      enum.UserType
}

type GetCorporationPostRequest struct {
	UserID        uint
	PostID        uint
	CorporationID uint
}

type GetPostRequest struct {
	UserID uint
	PostID uint
}
