package postgres

import "gorm.io/gorm"

type PaginationModifier struct {
	Limit  int
	Offset int
}

func NewPaginationModifier(limit, offset int) PaginationModifier {
	return PaginationModifier{
		Limit:  limit,
		Offset: offset,
	}
}

func (pagination PaginationModifier) Apply(query interface{}) interface{} {
	if db, ok := query.(*gorm.DB); ok {
		return db.Limit(pagination.Limit).Offset(pagination.Offset)
	}
	return query
}

type SortingModifier struct {
	Column string
	Desc   bool
}

func NewSortingModifier(column string, desc bool) SortingModifier {
	return SortingModifier{
		Column: column,
		Desc:   desc,
	}
}

func (sorting SortingModifier) Apply(query interface{}) interface{} {
	if db, ok := query.(*gorm.DB); ok {
		order := sorting.Column
		if sorting.Desc {
			order += " DESC"
		}
		return db.Order(order)
	}
	return query
}
