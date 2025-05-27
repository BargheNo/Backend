package blogdto

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/internal/domain/enum"
)

type CreatePostRequest struct {
	Title         string                `json:"title"`
	Content       string                `json:"content_html"`
	AuthorID      uint                  `json:"-"`
	CorporationID uint                  `json:"-"`
	CoverImage    *multipart.FileHeader `json:"cover_image"`
	Status        enum.PostStatus       `json:"status"`
}

type EditPostRequest struct {
	PostID        uint                  `json:"-"`
	AuthorID      uint                  `json:"-"`
	CorporationID uint                  `json:"-"`
	Title         *string               `json:"title"`
	Content       *string               `json:"content_html"`
	CoverImage    *multipart.FileHeader `json:"cover_image"`
	Status        uint                  `json:"status"`
}

type GetCorporationPostsRequest struct {
	CorporationID uint `uri:"corporationID" validate:"required"`
	Offset        int  `query:"offset" validate:"required"`
	Limit         int  `query:"limit" validate:"required"`
}

type DeletePostRequest struct {
	PostIDs       []uint `json:"-"`
	AuthorID      uint   `json:"-"`
	CorporationID uint   `json:"-"`
}

type AddPostMediaRequest struct {
	PostID        uint                  `json:"-"`
	AuthorID      uint                  `json:"-"`
	Media         *multipart.FileHeader `json:"media" validate:"required"`
	CorporationID uint                  `json:"-"`
}

type AccessPostMediaRequest struct {
	PostID        uint `json:"-"`
	AuthorID      uint `json:"-"`
	MediaID       uint `json:"-"`
	CorporationID uint `json:"-"`
}
