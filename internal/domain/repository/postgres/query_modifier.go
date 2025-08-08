package postgres

// QueryModifier defines how to modify queries
type QueryModifier interface {
	Apply(query interface{}) interface{}
}

// PaginationOptions defines pagination parameters
type PaginationOptions struct {
	Limit  int
	Offset int
}

// SortingOptions defines sorting parameters
type SortingOptions struct {
	Column string
	Asc    bool
}

// QueryOptions combines all query modification options
type QueryOptions struct {
	Pagination *PaginationOptions
	Sorting    *SortingOptions
}

// NewQueryOptions creates a new QueryOptions instance
func NewQueryOptions() *QueryOptions {
	return &QueryOptions{}
}

// WithPagination adds pagination to query options
func (q *QueryOptions) WithPagination(limit, offset int) *QueryOptions {
	q.Pagination = &PaginationOptions{
		Limit:  limit,
		Offset: offset,
	}
	return q
}

// WithSorting adds sorting to query options
func (q *QueryOptions) WithSorting(column string, asc bool) *QueryOptions {
	q.Sorting = &SortingOptions{
		Column: column,
		Asc:    asc,
	}
	return q
}

// HasPagination checks if pagination is set
func (q *QueryOptions) HasPagination() bool {
	return q.Pagination != nil
}

// HasSorting checks if sorting is set
func (q *QueryOptions) HasSorting() bool {
	return q.Sorting != nil
}
