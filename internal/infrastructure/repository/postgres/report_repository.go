package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
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

func (repo *ReportRepository) GetReportsByObjectType(db database.Database, objectType string, statuses []enum.ReportStatus, options *postgres.QueryOptions) ([]*entity.Report, error) {
	var reports []*entity.Report
	query := db.GetDB().Where("object_type = ? AND status IN ?", objectType, statuses)
	query = applyQueryOptions(query, options)

	result := query.Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}
	return reports, nil
}

func (repo *ReportRepository) CountReportsByObjectType(db database.Database, objectType string, statuses []enum.ReportStatus) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.Report{}).
		Where("object_type = ? AND status IN ?", objectType, statuses).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
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

func (repo *ReportRepository) FindMaintenanceReportsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.Report, error) {
	var reports []*entity.Report
	result := db.GetDB().
		Where("description ILIKE ? OR object_type = ?", "%"+query+"%", "maintenance")
	result = applyQueryOptions(result, options)
	result = result.Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}
	return reports, nil
}

func (repo *ReportRepository) CountMaintenanceReportsByQuery(db database.Database, query string) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.Report{}).
		Where("description ILIKE ? OR object_type = ?", "%"+query+"%", "maintenance").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ReportRepository) FindPanelReportsByQuery(db database.Database, query string, options *postgres.QueryOptions) ([]*entity.Report, error) {
	var reports []*entity.Report
	result := db.GetDB().
		Where("description ILIKE ? OR object_type = ?", "%"+query+"%", "panel")
	result = applyQueryOptions(result, options)
	result = result.Find(&reports)
	if result.Error != nil {
		return nil, result.Error
	}
	return reports, nil
}

func (repo *ReportRepository) CountPanelReportsByQuery(db database.Database, query string) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.Report{}).
		Where("description ILIKE ? OR object_type = ?", "%"+query+"%", "panel").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}
