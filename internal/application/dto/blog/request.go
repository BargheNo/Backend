package blogdto

import "mime/multipart"

type CreatePostRequest struct {
	Title         string                `json:"title"`
	Content       string                `json:"content_html"`
	AuthorID      uint                  `json:"-"`
	CorporationID uint                  `json:"-"`
	CoverImage    *multipart.FileHeader `json:"cover_image"`
}

type GetCorporationPostsRequest struct {
	CorporationID uint `uri:"corporationID" validate:"required"`
	Offset        int  `query:"offset" validate:"required"`
	Limit         int  `query:"limit" validate:"required"`
}
