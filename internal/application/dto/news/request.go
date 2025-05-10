package newsdto

import "github.com/BargheNo/Backend/internal/domain/enum"

type CreateNewsRequest struct {
	Title    string
	Content  string
	AuthorID uint
	Status   enum.NewsStatus
}
