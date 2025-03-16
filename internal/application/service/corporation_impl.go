package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationService struct {
	constants             *bootstrap.Constants
	JWTService            service.JWTService
	db                    database.Database
	CorporationRepository repository.CorporationRepository
}

func NewCorporationService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	corporationRepository repository.CorporationRepository,
) *CorporationService {
	return &CorporationService{
		constants:             constants,
		JWTService:            jwtService,
		db:                    db,
		CorporationRepository: corporationRepository,
	}
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) error {
	panic("not implemented")
}
