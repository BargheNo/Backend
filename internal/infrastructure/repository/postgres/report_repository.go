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

func (repo *ReportRepository) GetReportsByObjectType(db database.Database, objectType string, opts ...repository.QueryModifier) ([]*entity.Report, error) {
	var reports []*entity.Report
	query := db.GetDB().Where("object_type = ?", objectType)
	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}
	return reports, nil
}

func (repo *ReportRepository) FindReportByID(db database.Database, id uint) (*entity.Report, error) {
	var report entity.Report
	err := db.GetDB().Where("id = ?", id).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (repo *ReportRepository) UpdateReport(db database.Database, report *entity.Report) error {
	return db.GetDB().Save(report).Error
}
