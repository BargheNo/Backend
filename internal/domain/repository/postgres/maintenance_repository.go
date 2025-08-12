package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type MaintenanceRepository interface {
	CreateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error
	CreateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
	FindCorporationRequestByStatus(db database.Database, requestID, corporationID uint, allowedStatus []enum.MaintenanceRequestStatus) (*entity.MaintenanceRequest, error)
	FindCorporationRequestsByStatus(db database.Database, corporationID uint, allowedStatuses []enum.MaintenanceRequestStatus, options *QueryOptions) ([]*entity.MaintenanceRequest, error)
	CountCorporationRequestsByStatus(db database.Database, corporationID uint, allowedStatuses []enum.MaintenanceRequestStatus) (int64, error)
	FindRecordByID(db database.Database, recordID uint) (*entity.MaintenanceRecord, error)
	FindRecordByRequestID(db database.Database, requestID uint) (*entity.MaintenanceRecord, error)
	FindRequestByID(db database.Database, requestID uint) (*entity.MaintenanceRequest, error)
	FindRequestsByCustomerID(db database.Database, customerID uint, allowedStatus []enum.MaintenanceRequestStatus, options *QueryOptions) ([]*entity.MaintenanceRequest, error)
	FindRequestsByCustomerIDAndQuery(db database.Database, customerID uint, allowedStatus []enum.MaintenanceRequestStatus, query string, options *QueryOptions) ([]*entity.MaintenanceRequest, error)
	CountRequestsByCustomerIDAndQuery(db database.Database, customerID uint, allowedStatus []enum.MaintenanceRequestStatus, query string) (int64, error)
	CountRequestsByCustomerID(db database.Database, customerID uint, allowedStatus []enum.MaintenanceRequestStatus) (int64, error)
	FindRequestsByPanelIDAndStatus(db database.Database, panelID uint, allowedStatus []enum.MaintenanceRequestStatus, options *QueryOptions) ([]*entity.MaintenanceRequest, error)
	CountRequestsByPanelIDAndStatus(db database.Database, panelID uint, allowedStatus []enum.MaintenanceRequestStatus) (int64, error)
	UpdateMaintenanceRecord(db database.Database, maintenanceRecord *entity.MaintenanceRecord) error
	UpdateMaintenanceRequest(db database.Database, maintenanceRequest *entity.MaintenanceRequest) error
}
