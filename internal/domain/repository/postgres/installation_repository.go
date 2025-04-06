package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRepository interface {
	FindRequestByStatus(db database.Database, status []enum.InstallationRequestStatus, opts ...QueryModifier) []*entity.InstallationRequest
	FindRequestByID(db database.Database, requestID uint) (*entity.InstallationRequest, bool)
	FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, opts ...QueryModifier) []*entity.InstallationRequest
	FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, bool)
	CreateRequest(db database.Database, request *entity.InstallationRequest) error
	CreatePanel(db database.Database, panel *entity.Panel) error
	FindCorporationPanels(db database.Database, corporationID uint, opts ...QueryModifier) []*entity.Panel
	// FindOwnerPanels(db database.Database, ownerID uint, opts ...QueryModifier) []*entity.Panel
}
