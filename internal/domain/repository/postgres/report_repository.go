package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type ReportRepository interface {
	CreateReport(db database.Database, report *entity.Report) error
	GetReportsByObjectType(db database.Database, objectType string, statuses []enum.ReportStatus, options *QueryOptions) ([]*entity.Report, error)
	CountReportsByObjectType(db database.Database, objectType string, statuses []enum.ReportStatus) (int64, error)
	FindReportByID(db database.Database, id uint) (*entity.Report, error)
	UpdateReport(db database.Database, report *entity.Report) error
	FindMaintenanceReportsByQuery(db database.Database, query string, options *QueryOptions) ([]*entity.Report, error)
	CountMaintenanceReportsByQuery(db database.Database, query string) (int64, error)
	FindPanelReportsByQuery(db database.Database, query string, options *QueryOptions) ([]*entity.Report, error)
	CountPanelReportsByQuery(db database.Database, query string) (int64, error)
}
