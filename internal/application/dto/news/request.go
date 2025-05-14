package newsdto

import "github.com/BargheNo/Backend/internal/domain/enum"

type CreateNewsRequest struct {
	Title    string
	Content  string
	AuthorID uint
	Status   enum.NewsStatus
}

type EditNewsRequest struct {
	NewsID   uint
	AuthorID uint
	Title    *string
	Content  *string
	Status   uint
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

type GetNewsListRequest struct {
	Statuses []uint
	Offset   int
	Limit    int
}
