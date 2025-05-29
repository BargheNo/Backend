package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type GuaranteeService struct {
	constants           *bootstrap.Constants
	corporationService  service.CorporationService
	guaranteeRepository repository.GuaranteeRepository
	db                  database.Database
}

func NewGuaranteeService(
	constants *bootstrap.Constants,
	corporationService service.CorporationService,
	guaranteeRepository repository.GuaranteeRepository,
	db database.Database,
) *GuaranteeService {
	return &GuaranteeService{
		constants:           constants,
		corporationService:  corporationService,
		guaranteeRepository: guaranteeRepository,
		db:                  db,
	}
}

func (guaranteeService *GuaranteeService) ValidateGuaranteeOwnerShip(guaranteeID, corporationID uint) error {
	_, exist := guaranteeService.guaranteeRepository.FindCorporationGuarantee(guaranteeService.db, guaranteeID, corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.Guarantee}
		return notFoundError
	}
	return nil
}

func (guaranteeService *GuaranteeService) mapGuaranteeToResponse(guarantee *entity.Guarantee) guaranteedto.GuaranteeResponse {
	terms := guaranteeService.guaranteeRepository.FindGuaranteeTerms(guaranteeService.db, guarantee.ID)
	termsResponse := make([]guaranteedto.GuaranteeTermResponse, len(terms))
	for i, term := range terms {
		termsResponse[i] = guaranteedto.GuaranteeTermResponse{
			Title:       term.Title,
			Description: term.Description,
			Limitations: term.Limitations,
		}
	}

	return guaranteedto.GuaranteeResponse{
		ID:             guarantee.ID,
		Name:           guarantee.Name,
		Status:         guarantee.Status.String(),
		GuaranteeType:  guarantee.GuaranteeType.String(),
		DurationMonths: guarantee.DurationMonths,
		Description:    guarantee.Description,
		Terms:          termsResponse,
	}
}

func (guaranteeService *GuaranteeService) GetGuarantee(guaranteeID uint) (guaranteedto.GuaranteeResponse, error) {
	guarantee, exist := guaranteeService.guaranteeRepository.FindGuaranteeByID(guaranteeService.db, guaranteeID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.Guarantee}
		return guaranteedto.GuaranteeResponse{}, notFoundError
	}

	guaranteeDetails := guaranteeService.mapGuaranteeToResponse(guarantee)

	return guaranteeDetails, nil
}

func (guaranteeService *GuaranteeService) GetGuaranteeTypes() []guaranteedto.GuaranteeTypesResponse {
	types := enum.GetAllGuaranteeTypes()
	response := make([]guaranteedto.GuaranteeTypesResponse, len(types))
	for i, guaranteeType := range types {
		response[i] = guaranteedto.GuaranteeTypesResponse{
			ID:   uint(guaranteeType),
			Name: guaranteeType.String(),
		}
	}
	return response
}

func (guaranteeService *GuaranteeService) GetGuaranteeStatuses() []guaranteedto.GuaranteeTypesResponse {
	statuses := enum.GetAllGuaranteeStatuses()
	response := make([]guaranteedto.GuaranteeTypesResponse, len(statuses))
	for i, status := range statuses {
		response[i] = guaranteedto.GuaranteeTypesResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response
}

func (guaranteeService *GuaranteeService) GetCorporationGuarantee(request guaranteedto.GetGuaranteeRequest) guaranteedto.GuaranteeResponse {
	guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	guarantee, err := guaranteeService.GetGuarantee(request.GuaranteeID)
	if err != nil {
		panic(err)
	}
	return guarantee
}

func (guaranteeService *GuaranteeService) GetCorporationGuarantees(request guaranteedto.GetGuaranteesRequest) []guaranteedto.GuaranteeResponse {
	guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	allowedStatus := []enum.GuaranteeStatus{enum.GuaranteeStatus(request.Status)}
	if enum.GuaranteeStatus(request.Status) == enum.GuaranteeStatusAll {
		allowedStatus = enum.GetAllGuaranteeStatuses()
	}

	guarantees := guaranteeService.guaranteeRepository.FindCorporationGuarantees(guaranteeService.db, request.CorporationID, allowedStatus)
	response := make([]guaranteedto.GuaranteeResponse, len(guarantees))

	for i, guarantee := range guarantees {
		response[i] = guaranteeService.mapGuaranteeToResponse(guarantee)
	}
	return response
}

func (guaranteeService *GuaranteeService) AddGuarantee(request guaranteedto.CreateGuaranteeRequest) uint {
	guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	_, exist := guaranteeService.guaranteeRepository.FindCorporationGuaranteeByName(guaranteeService.db, request.CorporationID, request.Name)
	if exist {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(guaranteeService.constants.Field.Name, guaranteeService.constants.Tag.AlreadyExist)
		panic(conflictErrors)
	}

	guarantee := &entity.Guarantee{
		CorporationID:  request.CorporationID,
		Name:           request.Name,
		Status:         request.Status,
		GuaranteeType:  enum.GuaranteeType(request.GuaranteeType),
		DurationMonths: request.Duration,
		Description:    request.Description,
	}
	if err := guaranteeService.guaranteeRepository.CreateGuarantee(guaranteeService.db, guarantee); err != nil {
		panic(err)
	}
	for _, terms := range request.GuaranteeTermsRequest {
		if err := guaranteeService.addGuaranteeTerm(terms, guarantee.ID); err != nil {
			panic(err)
		}
	}
	return guarantee.ID
}

func (guaranteeService *GuaranteeService) addGuaranteeTerm(terms guaranteedto.GuaranteeTermsRequest, guaranteeID uint) error {
	guaranteeTerms := &entity.GuaranteeTerm{
		GuaranteeID: guaranteeID,
		Title:       terms.Title,
		Description: terms.Description,
		Limitations: terms.Limitations,
	}
	if err := guaranteeService.guaranteeRepository.CreateGuaranteeTerms(guaranteeService.db, guaranteeTerms); err != nil {
		return err
	}
	return nil
}

func (guaranteeService *GuaranteeService) UpdateGuaranteeStatus(request guaranteedto.ChangeStatusRequest) {
	guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)

	guarantee, exist := guaranteeService.guaranteeRepository.FindGuaranteeByID(guaranteeService.db, request.GuaranteeID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.Guarantee}
		panic(notFoundError)
	}

	if !enum.GuaranteeStatus(request.Status).IsValid() {
		return
	}

	if guarantee.Status == enum.GuaranteeStatus(request.Status) {
		var conflictErrors exception.ConflictErrors
		switch guarantee.Status {
		case enum.GuaranteeStatusActive:
			conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.AlreadyActive)
			panic(conflictErrors)
		case enum.GuaranteeStatusArchive:
			conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.AlreadyArchived)
			panic(conflictErrors)
		default:
			conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.StatusNotChange)
			panic(conflictErrors)
		}
	}

	guarantee.Status = enum.GuaranteeStatus(request.Status)

	if err := guaranteeService.guaranteeRepository.UpdateGuarantee(guaranteeService.db, guarantee); err != nil {
		panic(err)
	}
}
