package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func applyQueryOptions(query *gorm.DB, options *postgres.QueryOptions) *gorm.DB {
	if options == nil {
		return query
	}

	if options.Sorting != nil {
		if options.Sorting.Asc {
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: options.Sorting.Column},
			})
		} else {
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: options.Sorting.Column},
				Desc:   true,
			})
		}
	}

	if options.Pagination != nil {
		query = query.Offset(options.Pagination.Offset).Limit(options.Pagination.Limit)
	}

	return query
}
