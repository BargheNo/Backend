package blogdto

import "mime/multipart"

type CreatePostRequest struct {
	Title         string                `json:"title"`
	Content       string                `json:"content_html"`
	AuthorID      uint                  `json:"-"`
	CorporationID uint                  `json:"-"`
	CoverImage    *multipart.FileHeader `json:"cover_image"`
}
