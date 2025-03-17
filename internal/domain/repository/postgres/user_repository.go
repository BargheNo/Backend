package repository

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUserByID(db database.Database, id uint) (*entity.User, bool)
	FindUserByPhone(db database.Database, phone string) (*entity.User, bool)
	CreateUser(db database.Database, user *entity.User) error
	DeleteUserByPhone(db database.Database, phone string) error
	UpdateUser(db database.Database, user *entity.User) error
	FindUserRoles(db database.Database, user *entity.User) error
	FindRolePermissions(db database.Database, role *entity.Role) error
}
