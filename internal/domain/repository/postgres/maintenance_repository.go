package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRepository interface {
	CreateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error
	CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
	FindCorporationRequestByStatus(db database.Database, requestID, corporationID uint, allowedStatus []enum.MaintenanceRequestStatus) (*entity.MaintenanceRequest, bool)
	FindCorporationRequestsByStatus(db database.Database, corporationID uint, allowedStatuses []enum.MaintenanceRequestStatus, opts ...QueryModifier) []*entity.MaintenanceRequest
	FindRecordByID(db database.Database, recordID uint) (*entity.MaintenanceRecord, bool)
	FindRecordByRequestID(db database.Database, requestID uint) (*entity.MaintenanceRecord, bool)
	FindRequestByID(db database.Database, requestID uint) (*entity.MaintenanceRequest, bool)
	FindRequestsByCustomerID(db database.Database, customerID uint, allowedStatus []enum.MaintenanceRequestStatus, opts ...QueryModifier) []*entity.MaintenanceRequest
	FindRequestsByPanelIDAndStatus(db database.Database, panelID uint, allowedStatus []enum.MaintenanceRequestStatus, opts ...QueryModifier) []*entity.MaintenanceRequest
	UpdateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error
	UpdateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
}
