package service

import (
	"github.com/BargheNo/Backend/bootstrap"
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
