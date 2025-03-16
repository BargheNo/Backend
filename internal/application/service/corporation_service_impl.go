package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type CorporationService struct {
	constants             *bootstrap.Constants
	JWTService            service.JWTService
	db                    database.Database
	CorporationRepository repository.CorporationRepository
	CINService            service.CINService
}

func NewCorporationService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	corporationRepository repository.CorporationRepository,
	cinService service.CINService,
) *CorporationService {
	return &CorporationService{
		constants:             constants,
		JWTService:            jwtService,
		db:                    db,
		CorporationRepository: corporationRepository,
		CINService:            cinService,
	}
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) {
	var conflictErrors exception.ConflictErrors
	_, err := corporationService.CINService.ValidateCIN(registerInfo.CIN)
	if err != nil {
		panic(err)
	}
	_, corporationExists := corporationService.CorporationRepository.FindCorporationByCIN(corporationService.db, registerInfo.CIN)
	if corporationExists {
		conflictErrors.Add(corporationService.constants.Field.Corporation, corporationService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}

	corporation := &entity.Corporation{
		Name:     registerInfo.Name,
		CIN:      registerInfo.CIN,
		Password: registerInfo.Password,
		Status:   corporationService.constants.Status.Pending,
	}

	err = corporationService.CorporationRepository.CreateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}
