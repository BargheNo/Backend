package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MonitoringRepository struct{}

func NewMonitoringRepository() *MonitoringRepository {
	return &MonitoringRepository{}
}

func (repo *MonitoringRepository) FindPanelStatusByPanelID(db database.Database, panelID uint, options *postgres.QueryOptions) ([]*entity.PanelStatus, error) {
	var panelStatus []*entity.PanelStatus
	query := db.GetDB().Where("panel_id = ?", panelID)
	query = applyQueryOptions(query, options)

	result := query.Find(&panelStatus)
	if result.Error != nil {
		return nil, result.Error
	}
	return panelStatus, nil
}

func (repo *MonitoringRepository) CountPanelStatusByPanelID(db database.Database, panelID uint) (int64, error) {
	var count int64

	err := db.GetDB().
		Model(&entity.PanelStatus{}).
		Where("panel_id = ?", panelID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MonitoringRepository) FindPanelHistoryByPanelID(db database.Database, panelID uint, options *postgres.QueryOptions) ([]*entity.PanelHistory, error) {
	var panelHistory []*entity.PanelHistory
	query := db.GetDB().Where("panel_id = ?", panelID)
	query = applyQueryOptions(query, options)

	result := query.Find(&panelHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return panelHistory, nil
}

func (repo *MonitoringRepository) CountPanelHistoryByPanelID(db database.Database, panelID uint) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.PanelHistory{}).
		Where("panel_id = ?", panelID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MonitoringRepository) FindPanelEventByPanelID(db database.Database, panelID uint, options *postgres.QueryOptions) ([]*entity.PanelEvent, error) {
	var panelEvent []*entity.PanelEvent
	query := db.GetDB().Where("panel_id = ?", panelID)
	query = applyQueryOptions(query, options)

	result := query.Find(&panelEvent)
	if result.Error != nil {
		return nil, result.Error
	}
	return panelEvent, nil
}

func (repo *MonitoringRepository) CountPanelEventByPanelID(db database.Database, panelID uint) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.PanelEvent{}).
		Where("panel_id = ?", panelID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *MonitoringRepository) CreatePanelStatus(db database.Database, panelStatus *entity.PanelStatus) error {
	return db.GetDB().Create(panelStatus).Error
}

func (repo *MonitoringRepository) CreatePanelHistory(db database.Database, panelHistory *entity.PanelHistory) error {
	return db.GetDB().Create(panelHistory).Error
}

func (repo *MonitoringRepository) CreatePanelEvent(db database.Database, panelEvent *entity.PanelEvent) error {
	return db.GetDB().Create(panelEvent).Error
}
