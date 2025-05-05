package seed

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

type RoleSeeder struct {
	superAdmin             *bootstrap.AdminCredentials
	userRepository         repository.UserRepository
	notificationRepository repository.NotificationRepository
	db                     database.Database
}

func NewRoleSeeder(
	superAdmin *bootstrap.AdminCredentials,
	userRepository repository.UserRepository,
	notificationRepository repository.NotificationRepository,
	db database.Database,
) *RoleSeeder {
	return &RoleSeeder{
		superAdmin:             superAdmin,
		userRepository:         userRepository,
		notificationRepository: notificationRepository,
		db:                     db,
	}
}

func (roleSeeder *RoleSeeder) SeedRoles() {
	permissions := roleSeeder.seedPermissions()
	roles := roleSeeder.seedRolesWithPermissions(permissions)
	superAdminRoles := []*entity.Role{roles[enum.SuperAdmin]}
	roleSeeder.seedSuperAdmin(roleSeeder.superAdmin, superAdminRoles)
}

func (roleSeeder *RoleSeeder) seedPermissions() map[enum.PermissionType]*entity.Permission {
	permissions := make(map[enum.PermissionType]*entity.Permission)
	permissionTypes := enum.GetAllPermissionTypes()

	for _, permissionType := range permissionTypes {
		permission := roleSeeder.getOrCreatePermission(permissionType)
		permissions[permissionType] = permission
	}

	return permissions
}

func (roleSeeder *RoleSeeder) getOrCreatePermission(permissionType enum.PermissionType) *entity.Permission {
	permission, exists := roleSeeder.userRepository.FindPermissionByType(roleSeeder.db, permissionType)
	if exists {
		return permission
	}

	permission = &entity.Permission{
		Type:        permissionType,
		Description: permissionType.Description(),
		Category:    permissionType.Category(),
	}

	if err := roleSeeder.userRepository.CreatePermission(roleSeeder.db, permission); err != nil {
		panic(err)
	}

	return permission
}

func (roleSeeder *RoleSeeder) seedRolesWithPermissions(permissions map[enum.PermissionType]*entity.Permission) map[enum.RoleName]*entity.Role {
	roles := make(map[enum.RoleName]*entity.Role)
	roleNames := enum.GetAllRoleNames()
	for _, roleName := range roleNames {
		role := roleSeeder.getOrCreateRole(roleName)
		roles[roleName] = role
		roleSeeder.assignPermissionsToRole(role, roleName, permissions)
	}
	return roles
}

func (roleSeeder *RoleSeeder) getOrCreateRole(roleName enum.RoleName) *entity.Role {
	role, exists := roleSeeder.userRepository.FindRoleByName(roleSeeder.db, roleName.String())
	if exists {
		return role
	}

	role = &entity.Role{
		Name: roleName.String(),
	}

	if err := roleSeeder.userRepository.CreateRole(roleSeeder.db, role); err != nil {
		panic(err)
	}

	return role
}

func (roleSeeder *RoleSeeder) assignPermissionsToRole(role *entity.Role, roleName enum.RoleName, permissions map[enum.PermissionType]*entity.Permission) {
	for _, permissionType := range roleName.Permissions() {
		permission := permissions[permissionType]
		if !roleSeeder.userRepository.RoleHasPermission(roleSeeder.db, role.ID, permission.ID) {
			if err := roleSeeder.userRepository.AssignPermissionToRole(roleSeeder.db, role, permission); err != nil {
				panic(err)
			}
		}
	}
}

func (roleSeeder *RoleSeeder) seedSuperAdmin(adminCred *bootstrap.AdminCredentials, roles []*entity.Role) {
	admin := roleSeeder.getOrCreateAdmin(adminCred)
	for _, role := range roles {
		if !roleSeeder.userRepository.UserHasRole(roleSeeder.db, admin.ID, role.ID) {
			if err := roleSeeder.userRepository.AssignRoleToUser(roleSeeder.db, admin, role); err != nil {
				panic(err)
			}
		}

	}
}

func (roleSeeder *RoleSeeder) getOrCreateAdmin(adminCred *bootstrap.AdminCredentials) *entity.User {
	admin, exist := roleSeeder.userRepository.FindUserByPhone(roleSeeder.db, adminCred.Phone)
	if !exist {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminCred.Password), 14)
		if err != nil {
			panic(err)
		}
		admin = &entity.User{
			FirstName:     adminCred.FirstName,
			LastName:      adminCred.LastName,
			Phone:         adminCred.Phone,
			PhoneVerified: true,
			Password:      string(hashedPassword),
			Email:         adminCred.Email,
			EmailVerified: true,
			NationalCode:  adminCred.NationalCode,
			Status:        enum.UserStatusActive,
		}
		err = roleSeeder.userRepository.CreateUser(roleSeeder.db, admin)
		if err != nil {
			panic(err)
		}
	}

	notificationTypes := roleSeeder.notificationRepository.GetNotificationTypes(roleSeeder.db)
	for _, notificationType := range notificationTypes {
		_, exist := roleSeeder.notificationRepository.GetNotificationSettingByUserAndType(roleSeeder.db, admin.ID, notificationType.ID)
		if exist {
			continue
		}
		setting := &entity.NotificationSetting{
			UserID:         admin.ID,
			TypeID:         notificationType.ID,
			IsEmailEnabled: notificationType.SupportsEmail,
			IsPushEnabled:  notificationType.SupportsPush,
		}
		err := roleSeeder.notificationRepository.CreateNotificationSetting(roleSeeder.db, setting)
		if err != nil {
			panic(err)
		}
	}

	return admin
}
