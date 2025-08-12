package postgres

import (
	entity "github.com/BargheNo/Backend/internal/domain/entity/monitoring"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MonitoringRepository interface {
	FindPanelStatusByPanelID(db database.Database, panelID uint, options *QueryOptions) ([]*entity.PanelStatus, error)
	FindPanelHistoryByPanelID(db database.Database, panelID uint, options *QueryOptions) ([]*entity.PanelHistory, error)
	FindPanelEventByPanelID(db database.Database, panelID uint, options *QueryOptions) ([]*entity.PanelEvent, error)
}
