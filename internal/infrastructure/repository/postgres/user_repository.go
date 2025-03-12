package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) Create(db database.Database, user *entity.User) error {
	return db.GetDB().Create(user).Error
}

func (repo *UserRepository) Delete(db database.Database, userID uint) error {
	return db.GetDB().Delete(&entity.User{}, userID).Error
}
