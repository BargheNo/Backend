package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRepository interface {
	FindRequestByID(db database.Database, requestID uint) (*entity.InstallationRequest, bool)
	FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, opts ...QueryModifier) []*entity.InstallationRequest
	FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, bool)
	CreateRequest(db database.Database, request *entity.InstallationRequest) error
}
