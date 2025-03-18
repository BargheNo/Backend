package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRepository interface {
	FindUserByID(db database.Database, ownerID uint, status []enum.InstallationRequestStatus) []*entity.InstallationRequest
	CreateRequest(db database.Database, request *entity.InstallationRequest) error
}
