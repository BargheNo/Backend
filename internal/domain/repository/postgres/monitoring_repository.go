package postgres

import (
	entity "github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MonitoringRepository interface {
	FindPanelStatusByPanelID(db database.Database, panelID uint, options *QueryOptions) ([]*entity.PanelStatus, error)
	CountPanelStatusByPanelID(db database.Database, panelID uint) (int64, error)
	FindPanelHistoryByPanelID(db database.Database, panelID uint, options *QueryOptions) ([]*entity.PanelHistory, error)
	CountPanelHistoryByPanelID(db database.Database, panelID uint) (int64, error)
	FindPanelEventByPanelID(db database.Database, panelID uint, options *QueryOptions) ([]*entity.PanelEvent, error)
	CountPanelEventByPanelID(db database.Database, panelID uint) (int64, error)
	CreatePanelStatus(db database.Database, panelStatus *entity.PanelStatus) error
	CreatePanelHistory(db database.Database, panelHistory *entity.PanelHistory) error
	CreatePanelEvent(db database.Database, panelEvent *entity.PanelEvent) error
}
