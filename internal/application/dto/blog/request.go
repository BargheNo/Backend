package blogdto

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreatePostRequest struct {
	Title         string
	Content       string
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
	CoverImage    *multipart.FileHeader
	Status        uint
}

type GetPostsRequest struct {
	UserID        uint
	Statuses      []uint
	CorporationID uint
	Offset        int
	Limit         int
	UserType      enum.UserType
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
	AuthorID      uint
	MediaID       uint
	CorporationID uint
}
