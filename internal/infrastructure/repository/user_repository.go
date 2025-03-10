package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type SampleRepository struct{}

func NewSampleRepository() *SampleRepository {
	return &SampleRepository{}
}

func (repo *SampleRepository) Create(db database.Database, user *entity.User) error {
	return db.GetDB().Create(user).Error
}

func (repo *SampleRepository) Delete(db database.Database, userID uint) error {
	return db.GetDB().Delete(&entity.User{}, userID).Error
}
