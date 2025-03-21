package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type InstallationRepository struct{}

func NewInstallationRepository() *InstallationRepository {
	return &InstallationRepository{}
}

func (repo *InstallationRepository) FindOwnerRequests(db database.Database, ownerID uint, status []enum.InstallationRequestStatus) []*entity.InstallationRequest {
	var requests []*entity.InstallationRequest
	result := db.GetDB().Where("owner_id = ? and status IN ?", ownerID, status).Find(&requests)
	if result.Error != nil {
		panic(result.Error)
	}
	return requests
}

func (repo *InstallationRepository) FindOwnerRequestByName(db database.Database, ownerID uint, status []enum.InstallationRequestStatus, name string) (*entity.InstallationRequest, bool) {
	var request *entity.InstallationRequest
	result := db.GetDB().Where("owner_id = ? and status IN ? and name = ?", ownerID, status, name).Find(&request)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return request, true
}

func (repo *InstallationRepository) CreateRequest(db database.Database, request *entity.InstallationRequest) error {
	return db.GetDB().Create(&request).Error

}
