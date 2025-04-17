package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ReportRepository interface {
	CreateReport(db database.Database, report *entity.Report) error
	GetReports(db database.Database, opts ...QueryModifier) []*entity.Report
	GetReportByID(db database.Database, id uint) (*entity.Report, bool)
	UpdateReport(db database.Database, report *entity.Report) error
}
