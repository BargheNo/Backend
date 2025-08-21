package postgres

import (
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type RBACRepository struct{}

func NewRBACRepository() *RBACRepository {
	return &RBACRepository{}
}

func (repo *RBACRepository) FindAllPermissions(db database.Database, userType enum.UserType, options *postgres.QueryOptions) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	query := db.GetDB().Where("user_type = ?", userType).Find(&permissions)
	query = applyQueryOptions(query, options)

	result := query.Find(&permissions)
	if result.Error != nil {
		return nil, result.Error
	}
	return permissions, nil
}

func (repo *RBACRepository) CountAllPermissions(db database.Database, userType enum.UserType) (int64, error) {
	var count int64
	err := db.GetDB().Model(&entity.Permission{}).Where("user_type = ?", userType).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RBACRepository) FindRolePermissions(db database.Database, role *entity.Role) error {
	return db.GetDB().Preload("Permissions").First(&role).Error
}

func (repo *RBACRepository) FindAllRoles(db database.Database, userType enum.UserType, options *postgres.QueryOptions) ([]*entity.Role, error) {
	var roles []*entity.Role
	query := db.GetDB().Where("user_type = ?", userType).Find(&roles)
	query = applyQueryOptions(query, options)

	result := query.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (repo *RBACRepository) CountAllRoles(db database.Database, userType enum.UserType) (int64, error) {
	var count int64
	err := db.GetDB().Model(&entity.Role{}).Where("user_type = ?", userType).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RBACRepository) FindRolesByQuery(db database.Database, userType enum.UserType, query string, options *postgres.QueryOptions) ([]*entity.Role, error) {
	var roles []*entity.Role
	result := db.GetDB().Where("user_type = ? AND name ILIKE ?", userType, "%"+query+"%")
	result = applyQueryOptions(result, options)

	result = result.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (repo *RBACRepository) CountRolesByQuery(db database.Database, userType enum.UserType, query string) (int64, error) {
	var count int64
	err := db.GetDB().Model(&entity.Role{}).Where("user_type = ? AND name ILIKE ?", userType, "%"+query+"%").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RBACRepository) FindRoleByName(db database.Database, name string) (*entity.Role, error) {
	var role entity.Role
	result := db.GetDB().Where("name = ?", name).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

func (repo *RBACRepository) CreateRole(db database.Database, role *entity.Role) error {
	return db.GetDB().Create(&role).Error
}

func (repo *RBACRepository) FindPermissionByID(db database.Database, permissionID uint) (*entity.Permission, error) {
	var permission entity.Permission
	result := db.GetDB().First(&permission, permissionID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &permission, nil
}

func (repo *RBACRepository) FindRolesByPermission(db database.Database, permissionID uint, options *postgres.QueryOptions) ([]*entity.Role, error) {
	var roles []*entity.Role
	query := db.GetDB().
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Where("role_permissions.permission_id = ?", permissionID).
		Find(&roles)
	query = applyQueryOptions(query, options)

	result := query.Find(&roles)
	if result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (repo *RBACRepository) CountRolesByPermission(db database.Database, permissionID uint) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.Role{}).
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Where("role_permissions.permission_id = ?", permissionID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RBACRepository) AssignPermissionToRole(db database.Database, role *entity.Role, permission *entity.Permission) error {
	return db.GetDB().Model(role).Association("Permissions").Append(permission)
}

func (repo *RBACRepository) FindRoleByID(db database.Database, roleID uint) (*entity.Role, error) {
	var role entity.Role
	result := db.GetDB().First(&role, roleID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &role, nil
}

func (repo *RBACRepository) FindUsersByRoleID(db database.Database, roleID uint, options *postgres.QueryOptions) ([]*entity.User, error) {
	var users []*entity.User
	query := db.GetDB().
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users)
	query = applyQueryOptions(query, options)

	result := query.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *RBACRepository) CountUsersByRoleID(db database.Database, roleID uint) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RBACRepository) FindUsersByRoleIDAndQuery(db database.Database, roleID uint, query string, options *postgres.QueryOptions) ([]*entity.User, error) {
	var users []*entity.User
	result := db.GetDB().
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Find(&users)
	result = applyQueryOptions(result, options)

	result = result.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *RBACRepository) CountUsersByRoleIDAndQuery(db database.Database, roleID uint, query string) (int64, error) {
	var count int64
	err := db.GetDB().
		Model(&entity.User{}).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *RBACRepository) FindUserRoles(db database.Database, user *entity.User) error {
	return db.GetDB().Preload("Roles").First(&user).Error
}

func (repo *RBACRepository) DeleteRole(db database.Database, roleID uint) error {
	return db.GetDB().Unscoped().Delete(&entity.Role{}, roleID).Error
}

func (repo *RBACRepository) UpdateRole(db database.Database, role *entity.Role) error {
	return db.GetDB().Save(role).Error
}

func (repo *RBACRepository) ReplaceRolePermissions(db database.Database, role *entity.Role, permissions []*entity.Permission) error {
	return db.GetDB().Model(role).Association("Permissions").Replace(permissions)
}

func (repo *RBACRepository) ReplaceUserRoles(db database.Database, user *entity.User, roles []*entity.Role) error {
	return db.GetDB().Model(user).Association("Roles").Replace(roles)
}

func (repo *RBACRepository) ReplaceStaffRoles(db database.Database, staff *entity.CorporationStaff, roles []*entity.Role) error {
	return db.GetDB().Model(staff).Association("Roles").Replace(roles)
}

func (repo *RBACRepository) FindStaffRoles(db database.Database, staff *entity.CorporationStaff) error {
	return db.GetDB().Preload("Roles").First(staff).Error
}
