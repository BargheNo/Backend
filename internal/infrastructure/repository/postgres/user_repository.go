package repositoryimpl

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) FindUsers(db database.Database) ([]*entity.User, error) {
	var users []*entity.User
	result := db.GetDB().Find(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return users, nil
}

func (repo *UserRepository) FindUserByID(db database.Database, id uint) (*entity.User, bool) {
	var user entity.User
	result := db.GetDB().First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *UserRepository) FindUserByStatus(db database.Database, statuses []enum.UserStatus, opts ...repository.QueryModifier) []*entity.User {
	var users []*entity.User
	query := db.GetDB().Where("status IN ?", statuses)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&users)

	if result.Error != nil {
		panic(result.Error)
	}
	return users
}

func (repo *UserRepository) FindUserByEmail(db database.Database, email string) (*entity.User, bool) {
	var user entity.User
	result := db.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
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

func (repo *UserRepository) FindRoleByName(db database.Database, name string) (*entity.Role, bool) {
	var role entity.Role
	result := db.GetDB().Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &role, true
}

func (repo *UserRepository) CreateUser(db database.Database, user *entity.User) error {
	return db.GetDB().Create(&user).Error
}

func (repo *UserRepository) DeleteUserByPhone(db database.Database, phone string) error {
	return db.GetDB().Where("phone = ?", phone).Unscoped().Delete(&entity.User{}).Error
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

func (repo *UserRepository) FindPermissionByType(db database.Database, permissionType enum.PermissionType) (*entity.Permission, bool) {
	var permission entity.Permission
	result := db.GetDB().Where("type = ?", permissionType).First(&permission)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &permission, true
}

func (repo *UserRepository) RoleHasPermission(db database.Database, roleID uint, permissionID uint) bool {
	var count int64
	db.GetDB().
		Table("role_permissions").
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count)
	return count > 0
}

func (repo *UserRepository) UserHasRole(db database.Database, userID uint, roleID uint) bool {
	var count int64
	db.GetDB().
		Table("user_roles").
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count)
	return count > 0
}

func (repo *UserRepository) CreateRole(db database.Database, role *entity.Role) error {
	return db.GetDB().Create(&role).Error
}

func (repo *UserRepository) CreatePermission(db database.Database, permission *entity.Permission) error {
	return db.GetDB().Create(&permission).Error
}

func (repo *UserRepository) AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error {
	return db.GetDB().Model(role).Association("Permissions").Append(permission)
}

func (repo *UserRepository) AssignRoleToUser(db database.Database, user *entity.User, role *entity.Role) error {
	return db.GetDB().Model(user).Association("Roles").Append(role)
}

func (repo *UserRepository) FindAllPermissions(db database.Database) []*entity.Permission {
	var permissions []*entity.Permission
	result := db.GetDB().Find(&permissions)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return permissions
}

func (repo *UserRepository) FindAllRoles(db database.Database) []*entity.Role {
	var roles []*entity.Role
	result := db.GetDB().Find(&roles)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return roles
}

func (repo *UserRepository) FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, bool) {
	var permission entity.Permission
	result := db.GetDB().First(&permission, permissionID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &permission, true
}

func (repo *UserRepository) FindRoleByID(db database.Database, roleID uint) (*entity.Role, bool) {
	var role entity.Role
	result := db.GetDB().First(&role, roleID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &role, true
}

func (repo *UserRepository) FindUsersByRoleID(db database.Database, roleID uint) []*entity.User {
	var users []*entity.User
	result := db.GetDB().
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return users
}

func (repo *UserRepository) FindUsersByPermission(db database.Database, permissionTypes []enum.PermissionType) []*entity.User {
	var users []*entity.User

	result := db.GetDB().
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("permissions.type IN ?", permissionTypes).
		Distinct().
		Find(&users)

	if result.Error != nil {
		panic(result.Error)
	}

	return users
}

func (repo *UserRepository) DeleteRole(db database.Database, roleID uint) error {
	return db.GetDB().Unscoped().Delete(&entity.Role{}, roleID).Error
}

func (repo *UserRepository) UpdateRole(db database.Database, role *entity.Role) error {
	return db.GetDB().Save(&role).Error
}

func (repo *UserRepository) ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []entity.Permission) error {
	return db.GetDB().Model(&role).Association("Permissions").Replace(permissions)

}

func (repo *UserRepository) ReplaceUserRoles(db database.Database, user *entity.User, roles []entity.Role) error {
	return db.GetDB().Model(&user).Association("Roles").Replace(roles)

}
