package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) FindUserByPhone(db database.Database, phone string) (*entity.User, bool) {
	var user entity.User
	result := db.GetDB().Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *UserRepository) CreateUser(db database.Database, user *entity.User) error {
	return db.GetDB().Create(&user).Error
}

func (repo *UserRepository) DeleteUserByPhone(db database.Database, phone string) error {
	return db.GetDB().Where("phone = ?", phone).Delete(&entity.User{}).Error
}

func (repo *UserRepository) UpdateUser(db database.Database, user *entity.User) error {
	return db.GetDB().Save(&user).Error
}

func (repo *UserRepository) FindUserRoles(db database.Database, user *entity.User) error {
	return db.GetDB().Preload("Roles").First(&user).Error
}

func (repo *UserRepository) FindRolePermissions(db database.Database, role *entity.Role) error {
	return db.GetDB().Preload("Permissions").First(&role).Error
}
