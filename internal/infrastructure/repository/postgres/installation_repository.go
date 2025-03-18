package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationRepository struct{}

func NewInstallationRepository() *InstallationRepository {
	return &InstallationRepository{}
}

func (repo *InstallationRepository) FindUserByID(db database.Database, ownerID uint, status []enum.InstallationRequestStatus) []*entity.InstallationRequest {
	var requests []*entity.InstallationRequest
	result := db.GetDB().Where("owner_id = ? and status IN ?", ownerID, status).Find(&requests)
	if result.Error != nil {
		panic(result.Error)
	}
	return requests
}

func (repo *InstallationRepository) CreateRequest(db database.Database, request *entity.InstallationRequest) error {
	return db.GetDB().Create(&request).Error

}
