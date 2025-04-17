package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ReportRepository struct {
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

func (r *ReportRepository) CreateReport(db database.Database, report *entity.Report) error {
	return db.GetDB().Create(report).Error
}
