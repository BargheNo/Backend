package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ReportRepository struct {
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

func (r *ReportRepository) CreateReport(db database.Database, report *entity.Report) error {
	return db.GetDB().Create(report).Error
}

func (repo *ReportRepository) GetReports(db database.Database, opts ...repository.QueryModifier) []*entity.Report {
	var reports []*entity.Report
	query := db.GetDB().Model(&entity.Report{})

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	err := query.Find(&reports).Error
	if err != nil {
		return nil
	}
	return reports
}
