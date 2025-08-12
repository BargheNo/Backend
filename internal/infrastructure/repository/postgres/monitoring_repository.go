package postgres

import (
	entity "github.com/BargheNo/Backend/internal/domain/entity/monitoring"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MonitoringRepository struct{}

func NewMonitoringRepository() *MonitoringRepository {
	return &MonitoringRepository{}
}

func (r *MonitoringRepository) FindPanelStatusByPanelID(db database.Database, panelID uint, options *postgres.QueryOptions) ([]*entity.PanelStatus, error) {
	var panelStatus []*entity.PanelStatus
	query := db.GetDB().Where("panel_id = ?", panelID)
	query = applyQueryOptions(query, options)

	result := query.Find(&panelStatus)
	if result.Error != nil {
		return nil, result.Error
	}
	return panelStatus, nil
}

func (r *MonitoringRepository) FindPanelHistoryByPanelID(db database.Database, panelID uint, options *postgres.QueryOptions) ([]*entity.PanelHistory, error) {
	var panelHistory []*entity.PanelHistory
	query := db.GetDB().Where("panel_id = ?", panelID)
	query = applyQueryOptions(query, options)

	result := query.Find(&panelHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return panelHistory, nil
}

func (r *MonitoringRepository) FindPanelEventByPanelID(db database.Database, panelID uint, options *postgres.QueryOptions) ([]*entity.PanelEvent, error) {
	var panelEvent []*entity.PanelEvent
	query := db.GetDB().Where("panel_id = ?", panelID)
	query = applyQueryOptions(query, options)

	result := query.Find(&panelEvent)
	if result.Error != nil {
		return nil, result.Error
	}
	return panelEvent, nil
}
