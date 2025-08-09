package controller

func GetOffsetLimit(page, pageSize, defaultPage, defaultPageSize int) (int, int) {
	if page <= 0 {
		page = defaultPage
	}
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	return (page - 1) * pageSize, pageSize
}

type PaginationMeta struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalItems  int64 `json:"totalItems"`
	TotalPages  int   `json:"totalPages"`
	HasNextPage bool  `json:"hasNextPage"`
	HasPrevPage bool  `json:"hasPrevPage"`
}

type PaginatedResponse[T any] struct {
	Data       []T             `json:"data"`
	Pagination *PaginationMeta `json:"pagination"`
}

func NewPaginatedResponse[T any](data []T, totalItems int64, offset, limit int) *PaginatedResponse[T] {
	pageSize := limit
	currentPage := int(offset/limit) + 1

	totalPages := 0
	if pageSize > 0 {
		totalPages = int((totalItems + int64(pageSize) - 1) / int64(pageSize))
	}

	return &PaginatedResponse[T]{
		Data: data,
		Pagination: &PaginationMeta{
			CurrentPage: currentPage,
			PageSize:    pageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			HasNextPage: currentPage < totalPages,
			HasPrevPage: currentPage > 1,
		},
	}
}
