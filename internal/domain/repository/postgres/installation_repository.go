package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRepository interface {
	FindRequestsByStatus(db database.Database, status []enum.InstallationRequestStatus, options *QueryOptions) ([]*entity.InstallationRequest, error)
	CountRequestsByStatus(db database.Database, status []enum.InstallationRequestStatus) (int64, error)
	FindRequestsByQuery(db database.Database, query string, options *QueryOptions) ([]*entity.InstallationRequest, error)
	CountRequestsByQuery(db database.Database, query string) (int64, error)
	FindRequestByID(db database.Database, requestID uint) (*entity.InstallationRequest, error)
	FindRequestByOwner(db database.Database, requestID, ownerID uint) (*entity.InstallationRequest, error)
	FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, options *QueryOptions) ([]*entity.InstallationRequest, error)
	CountOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus) (int64, error)
	FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, error)
	FindCorporationPanel(db database.Database, panelID, corporationID uint) (*entity.Panel, error)
	FindCorporationPanels(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus, options *QueryOptions) ([]*entity.Panel, error)
	CountCorporationPanels(db database.Database, corporationID uint, allowedStatus []enum.PanelStatus) (int64, error)
	FindCustomerPanels(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, options *QueryOptions) ([]*entity.Panel, error)
	CountCustomerPanels(db database.Database, customerID uint, allowedStatus []enum.PanelStatus) (int64, error)
	FindCustomerPanelsByQuery(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, query string, options *QueryOptions) ([]*entity.Panel, error)
	CountCustomerPanelsByQuery(db database.Database, customerID uint, allowedStatus []enum.PanelStatus, query string) (int64, error)
	FindCustomerPanel(db database.Database, panelID, customerID uint) (*entity.Panel, error)
	CreateRequest(db database.Database, request *entity.InstallationRequest) error
	UpdateRequest(db database.Database, request *entity.InstallationRequest) error
	DeleteRequest(db database.Database, request *entity.InstallationRequest) error
	FindPanelByNameAndCustomerID(db database.Database, panelName string, customerID uint) (*entity.Panel, error)
	FindPanelByOwner(db database.Database, panelID, customerID uint) (*entity.Panel, error)
	FindPanelByID(db database.Database, panelID uint) (*entity.Panel, error)
	FindPanelsByStatus(db database.Database, allowedStatus []enum.PanelStatus, options *QueryOptions) ([]*entity.Panel, error)
	CountPanelsByStatus(db database.Database, allowedStatus []enum.PanelStatus) (int64, error)
	FindPanelsByQuery(db database.Database, query string, options *QueryOptions) ([]*entity.Panel, error)
	CountPanelsByQuery(db database.Database, query string) (int64, error)
	CreatePanel(db database.Database, panel *entity.Panel) error
	UpdatePanel(db database.Database, panel *entity.Panel) error
	DeletePanel(db database.Database, panel *entity.Panel) error
}
