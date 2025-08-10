package postgres

type QueryModifier interface {
	Apply(query interface{}) interface{}
}

type PaginationOptions struct {
	Limit  int
	Offset int
}

type SortingOptions struct {
	Column string
	Asc    bool
}

type QueryOptions struct {
	Pagination *PaginationOptions
	Sorting    *SortingOptions
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

func (q *QueryOptions) WithPagination(limit, offset int) *QueryOptions {
	q.Pagination = &PaginationOptions{
		Limit:  limit,
		Offset: offset,
	}
	return q
}

func (q *QueryOptions) WithSorting(column string, asc bool) *QueryOptions {
	q.Sorting = &SortingOptions{
		Column: column,
		Asc:    asc,
	}
	return q
}

func (q *QueryOptions) HasPagination() bool {
	return q.Pagination != nil
}

func (q *QueryOptions) HasSorting() bool {
	return q.Sorting != nil
}
