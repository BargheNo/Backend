package service

import (
	"github.com/BargheNo/Backend/bootstrap"
	rbacdto "github.com/BargheNo/Backend/internal/application/dto/rbac"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/enum/sortby"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type RBACService struct {
	constants      *bootstrap.Constants
	rbacRepository postgres.RBACRepository
	db             database.Database
}

func NewRBACService(
	constants *bootstrap.Constants,
	rbacRepository postgres.RBACRepository,
	db database.Database,
) *RBACService {
	return &RBACService{
		constants:      constants,
		rbacRepository: rbacRepository,
		db:             db,
	}
}

func (rbacService *RBACService) getSortByColumn(requested uint) string {
	allowed := sortby.GetUserSortableColumns()
	sortBy := sortby.UserSortBy(requested)
	if _, ok := allowed[sortBy]; ok {
		return sortBy.DBColumn()
	}
	return sortby.NewsSortByCreatedAt.DBColumn()
}

func (rbacService *RBACService) getRBACUserType(isStaff bool) enum.UserType {
	permissionUserType := enum.UserTypeAdmin
	if isStaff {
		permissionUserType = enum.UserTypeCorporation
	}
	return permissionUserType
}

// TODO: add some search query to this later
func (rbacService *RBACService) GetPermissions(request rbacdto.GetPermissionsListRequest) ([]rbacdto.PermissionResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	permissionUserType := rbacService.getRBACUserType(request.IsStaff)
	permissions, err := rbacService.rbacRepository.FindAllPermissions(rbacService.db, permissionUserType, options)
	if err != nil {
		return nil, 0, err
	}
	permissionsResponse := make([]rbacdto.PermissionResponse, len(permissions))

	for i, permission := range permissions {
		isCorpStaff := permission.UserType == enum.UserTypeCorporation
		permissionsResponse[i] = rbacdto.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Type.String(),
			Description: permission.Type.Description(),
			IsCorpStaff: isCorpStaff,
			Category:    permission.Category.String(),
		}
	}

	count, err := rbacService.rbacRepository.CountAllPermissions(rbacService.db, permissionUserType)
	if err != nil {
		return nil, 0, err
	}
	return permissionsResponse, count, nil
}

func (rbacService *RBACService) GetPermissionRoles(request rbacdto.GetPermissionRolesRequest) ([]rbacdto.RoleResponse, int64, error) {
	permission, err := rbacService.rbacRepository.FindPermissionByID(rbacService.db, request.PermissionID)
	if err != nil {
		return nil, 0, err
	}
	if permission == nil {
		notFoundError := exception.NotFoundError{Item: rbacService.constants.Field.Permission}
		return nil, 0, notFoundError
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(rbacService.getSortByColumn(request.SortBy), request.Asc)

	roles, err := rbacService.rbacRepository.FindRolesByPermission(rbacService.db, request.PermissionID, options)
	if err != nil {
		return nil, 0, err
	}
	rolesResponse := make([]rbacdto.RoleResponse, len(roles))
	for i, role := range roles {
		permissions, err := rbacService.getRolePermissions(role)
		if err != nil {
			return nil, 0, err
		}
		isCorpStaff := role.UserType == enum.UserTypeCorporation
		rolesResponse[i] = rbacdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			IsCorpStaff: isCorpStaff,
			Permissions: permissions,
		}
	}
	count, err := rbacService.rbacRepository.CountRolesByPermission(rbacService.db, request.PermissionID)
	if err != nil {
		return nil, 0, err
	}
	return rolesResponse, count, nil
}

func (rbacService *RBACService) GetUserPermissions(user *entity.User) ([]rbacdto.PermissionResponse, error) {
	var permissions []rbacdto.PermissionResponse
	if err := rbacService.rbacRepository.FindUserRoles(rbacService.db, user); err != nil {
		return nil, err
	}
	for _, role := range user.Roles {
		rolePermissions, err := rbacService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, rolePermissions...)
	}
	return permissions, nil
}

func (rbacService *RBACService) getRolePermissions(role *entity.Role) ([]rbacdto.PermissionResponse, error) {
	if err := rbacService.rbacRepository.FindRolePermissions(rbacService.db, role); err != nil {
		return nil, err
	}
	permissions := make([]rbacdto.PermissionResponse, len(role.Permissions))
	for i, permission := range role.Permissions {
		isCorpStaff := permission.UserType == enum.UserTypeCorporation
		permissions[i] = rbacdto.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Type.String(),
			Description: permission.Type.Description(),
			IsCorpStaff: isCorpStaff,
			Category:    permission.Category.String(),
		}
	}
	return permissions, nil
}

func (rbacService *RBACService) getRolesByQuery(query string, userType enum.UserType, options *postgres.QueryOptions) ([]*entity.Role, int64, error) {
	if query == "" {
		roles, err := rbacService.rbacRepository.FindAllRoles(rbacService.db, userType, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := rbacService.rbacRepository.CountAllRoles(rbacService.db, userType)
		if err != nil {
			return nil, 0, err
		}
		return roles, count, nil
	}
	roles, err := rbacService.rbacRepository.FindRolesByQuery(rbacService.db, userType, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := rbacService.rbacRepository.CountRolesByQuery(rbacService.db, userType, query)
	if err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}

func (rbacService *RBACService) GetRoles(request rbacdto.GetRolesListRequest) ([]rbacdto.RoleResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	roleUserType := rbacService.getRBACUserType(request.IsStaff)
	roles, count, err := rbacService.getRolesByQuery(request.Query, roleUserType, options)
	if err != nil {
		return nil, 0, err
	}
	rolesResponse := make([]rbacdto.RoleResponse, len(roles))
	for i, role := range roles {
		permissions, err := rbacService.getRolePermissions(role)
		if err != nil {
			return nil, 0, err
		}
		isCorpStaff := role.UserType == enum.UserTypeCorporation
		rolesResponse[i] = rbacdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			IsCorpStaff: isCorpStaff,
			Permissions: permissions,
		}
	}

	return rolesResponse, count, nil
}

func (rbacService *RBACService) checkDuplicateRole(roleName string) error {
	existingRole, err := rbacService.rbacRepository.FindRoleByName(rbacService.db, roleName)
	if err != nil {
		return err
	}
	if existingRole != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(rbacService.constants.Field.Role, rbacService.constants.Tag.AlreadyExist)
		return conflictErrors
	}
	return nil
}

func (rbacService *RBACService) getPermission(permissionID uint) (*entity.Permission, error) {
	permission, err := rbacService.rbacRepository.FindPermissionByID(rbacService.db, permissionID)
	if err != nil {
		return nil, err
	}

	if permission == nil {
		notFoundError := exception.NotFoundError{Item: rbacService.constants.Field.Permission}
		return nil, notFoundError
	}

	return permission, nil
}

func (rbacService *RBACService) assignPermissionsToRole(role *entity.Role, permissionIDs []uint, userType enum.UserType) error {
	err := rbacService.db.WithTransaction(func(tx database.Database) error {
		existingPermissions := make(map[uint]bool)
		for _, permissionID := range permissionIDs {
			if existingPermissions[permissionID] {
				continue
			}

			permission, err := rbacService.getPermission(permissionID)
			if err != nil {
				return err
			}

			if permission.UserType != userType {
				continue
			}

			if err := rbacService.rbacRepository.AssignPermissionToRole(rbacService.db, role, permission); err != nil {
				return err
			}
			existingPermissions[permissionID] = true
		}
		return nil
	})
	return err
}

func (rbacService *RBACService) CreateRole(newRoleRequest rbacdto.NewRoleRequest) error {
	if err := rbacService.checkDuplicateRole(newRoleRequest.Name); err != nil {
		return err
	}

	userType := rbacService.getRBACUserType(newRoleRequest.IsStaff)
	role := &entity.Role{
		Name:     newRoleRequest.Name,
		UserType: userType,
	}

	err := rbacService.db.WithTransaction(func(tx database.Database) error {
		err := rbacService.rbacRepository.CreateRole(tx, role)
		if err != nil {
			return err
		}

		if err := rbacService.assignPermissionsToRole(role, newRoleRequest.PermissionIDs, userType); err != nil {
			return err
		}

		return nil
	})
	return err
}

func (rbacService *RBACService) getRole(roleID uint) (*entity.Role, error) {
	role, err := rbacService.rbacRepository.FindRoleByID(rbacService.db, roleID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		notFoundError := exception.NotFoundError{Item: rbacService.constants.Field.Role}
		return nil, notFoundError
	}

	return role, nil
}

func (rbacService *RBACService) GetRoleDetails(roleID uint) (rbacdto.RoleResponse, error) {
	role, err := rbacService.getRole(roleID)
	if err != nil {
		return rbacdto.RoleResponse{}, err
	}

	permissions, err := rbacService.getRolePermissions(role)
	if err != nil {
		return rbacdto.RoleResponse{}, err
	}

	isCorpStaff := role.UserType == enum.UserTypeCorporation
	return rbacdto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		IsCorpStaff: isCorpStaff,
		Permissions: permissions,
	}, nil
}

func (rbacService *RBACService) getRoleOwnersByQuery(roleID uint, query string, options *postgres.QueryOptions) ([]*entity.User, int64, error) {
	if query == "" {
		users, err := rbacService.rbacRepository.FindUsersByRoleID(rbacService.db, roleID, options)
		if err != nil {
			return nil, 0, err
		}
		count, err := rbacService.rbacRepository.CountUsersByRoleID(rbacService.db, roleID)
		if err != nil {
			return nil, 0, err
		}
		return users, count, nil
	}

	users, err := rbacService.rbacRepository.FindUsersByRoleIDAndQuery(rbacService.db, roleID, query, options)
	if err != nil {
		return nil, 0, err
	}
	count, err := rbacService.rbacRepository.CountUsersByRoleIDAndQuery(rbacService.db, roleID, query)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (rbacService *RBACService) GetRoleOwners(request rbacdto.GetRoleOwnersRequest) ([]*entity.User, int64, error) {
	_, err := rbacService.getRole(request.RoleID)
	if err != nil {
		return nil, 0, err
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset).
		WithSorting(rbacService.getSortByColumn(request.SortBy), request.Asc)

	users, count, err := rbacService.getRoleOwnersByQuery(request.RoleID, request.Query, options)
	if err != nil {
		return nil, 0, err
	}
	return users, count, err
}

func (rbacService *RBACService) GetUserRoles(user *entity.User) ([]rbacdto.RoleResponse, error) {
	if err := rbacService.rbacRepository.FindUserRoles(rbacService.db, user); err != nil {
		return nil, err
	}
	roles := make([]rbacdto.RoleResponse, len(user.Roles))

	for i, role := range user.Roles {
		permissions, err := rbacService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		isCorpStaff := role.UserType == enum.UserTypeCorporation
		roles[i] = rbacdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			IsCorpStaff: isCorpStaff,
			Permissions: permissions,
		}
	}
	return roles, nil
}

func (rbacService *RBACService) GetStaffRoles(staff *entity.CorporationStaff) ([]rbacdto.RoleResponse, error) {
	if err := rbacService.rbacRepository.FindStaffRoles(rbacService.db, staff); err != nil {
		return nil, err
	}
	roles := make([]rbacdto.RoleResponse, len(staff.Roles))

	for i, role := range staff.Roles {
		permissions, err := rbacService.getRolePermissions(&role)
		if err != nil {
			return nil, err
		}
		isCorpStaff := role.UserType == enum.UserTypeCorporation
		roles[i] = rbacdto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			IsCorpStaff: isCorpStaff,
			Permissions: permissions,
		}
	}
	return roles, nil
}

func (rbacService *RBACService) DeleteRole(roleID uint) error {
	if _, err := rbacService.getRole(roleID); err != nil {
		return err
	}

	if err := rbacService.rbacRepository.DeleteRole(rbacService.db, roleID); err != nil {
		return err
	}
	return nil
}

func (rbacService *RBACService) UpdateRole(newRoleRequest rbacdto.UpdateRoleRequest) error {
	role, err := rbacService.getRole(newRoleRequest.RoleID)
	if err != nil {
		return err
	}

	existingPermissions := make(map[uint]bool)
	var permissions []*entity.Permission
	for _, permissionID := range newRoleRequest.PermissionIDs {
		if existingPermissions[permissionID] {
			continue
		}

		permission, err := rbacService.getPermission(permissionID)
		if err != nil {
			return err
		}
		if permission.UserType != role.UserType {
			continue
		}

		permissions = append(permissions, permission)
		existingPermissions[permissionID] = true
	}

	err = rbacService.db.WithTransaction(func(tx database.Database) error {
		if newRoleRequest.Name != nil {
			role.Name = *newRoleRequest.Name
			if err := rbacService.rbacRepository.UpdateRole(tx, role); err != nil {
				return err
			}
		}

		if err := rbacService.rbacRepository.ReplaceRolePermissions(tx, role, permissions); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (rbacService *RBACService) getUniqueRoles(roleIDs []uint) ([]*entity.Role, error) {
	existingRoles := make(map[uint]bool)
	var roles []*entity.Role
	for _, roleID := range roleIDs {
		if existingRoles[roleID] {
			continue
		}

		role, err := rbacService.getRole(roleID)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
		existingRoles[roleID] = true
	}
	return roles, nil
}

func (rbacService *RBACService) UpdateUserRoles(userRolesRequest rbacdto.UpdateUserRolesRequest) error {
	roles, err := rbacService.getUniqueRoles(userRolesRequest.RoleIDs)
	if err != nil {
		return err
	}

	if err := rbacService.rbacRepository.ReplaceUserRoles(rbacService.db, userRolesRequest.User, roles); err != nil {
		return err
	}

	return nil
}

func (rbacService *RBACService) UpdateStaffRoles(request rbacdto.UpdateStaffRolesRequest) error {
	roles, err := rbacService.getUniqueRoles(request.RoleIDs)
	if err != nil {
		return err
	}

	if err := rbacService.rbacRepository.ReplaceStaffRoles(rbacService.db, request.Staff, roles); err != nil {
		return err
	}

	return nil
}
