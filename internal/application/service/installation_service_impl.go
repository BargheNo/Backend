package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type InstallationService struct {
	constants *bootstrap.Constants
	db        database.Database
}

func NewInstallationService(
	constants *bootstrap.Constants,
	db database.Database,
) *InstallationService {
	return &InstallationService{
		constants: constants,
		db:        db,
	}
}
