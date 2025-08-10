package service

import (
	"github.com/BargheNo/Backend/bootstrap"
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type GuaranteeService struct {
	constants           *bootstrap.Constants
	corporationService  usecase.CorporationService
	userService         usecase.UserService
	guaranteeRepository postgres.GuaranteeRepository
	db                  database.Database
}

func NewGuaranteeService(
	constants *bootstrap.Constants,
	corporationService usecase.CorporationService,
	userService usecase.UserService,
	guaranteeRepository postgres.GuaranteeRepository,
	db database.Database,
) *GuaranteeService {
	return &GuaranteeService{
		constants:           constants,
		corporationService:  corporationService,
		userService:         userService,
		guaranteeRepository: guaranteeRepository,
		db:                  db,
	}
}

func (guaranteeService *GuaranteeService) ValidateActiveGuaranteeOwnerShip(guaranteeID, corporationID uint) error {
	guarantee, err := guaranteeService.guaranteeRepository.FindCorporationGuarantee(guaranteeService.db, guaranteeID, corporationID)
	if err != nil {
		return err
	}
	if guarantee == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.Guarantee}
		return notFoundError
	}

	if guarantee.Status != enum.GuaranteeStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.NotActive)
		return conflictErrors
	}

	return nil
}

func (guaranteeService *GuaranteeService) mapGuaranteeToResponse(guarantee *entity.Guarantee) (guaranteedto.GuaranteeResponse, error) {
	terms, err := guaranteeService.guaranteeRepository.FindGuaranteeTerms(guaranteeService.db, guarantee.ID)
	if err != nil {
		return guaranteedto.GuaranteeResponse{}, err
	}
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
	}, nil
}

func (guaranteeService *GuaranteeService) GetGuarantee(guaranteeID uint) (guaranteedto.GuaranteeResponse, error) {
	guarantee, err := guaranteeService.guaranteeRepository.FindGuaranteeByID(guaranteeService.db, guaranteeID)
	if err != nil {
		return guaranteedto.GuaranteeResponse{}, err
	}
	if guarantee == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.Guarantee}
		return guaranteedto.GuaranteeResponse{}, notFoundError
	}

	guaranteeDetails, err := guaranteeService.mapGuaranteeToResponse(guarantee)
	if err != nil {
		return guaranteedto.GuaranteeResponse{}, err
	}

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

func (guaranteeService *GuaranteeService) GetCorporationGuarantee(request guaranteedto.GetGuaranteeRequest) (guaranteedto.GuaranteeResponse, error) {
	err := guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)
	if err != nil {
		return guaranteedto.GuaranteeResponse{}, err
	}

	guarantee, err := guaranteeService.GetGuarantee(request.GuaranteeID)
	if err != nil {
		return guaranteedto.GuaranteeResponse{}, err
	}
	return guarantee, nil
}

func (guaranteeService *GuaranteeService) mapToFilterStatuses(enumStatus uint) []enum.GuaranteeStatus {
	statuses := enum.GetAllGuaranteeStatuses()
	for _, status := range statuses {
		if uint(status) == enumStatus {
			if status == enum.GuaranteeStatusAll {
				return statuses
			}
			return []enum.GuaranteeStatus{status}
		}
	}
	return statuses
}

func (guaranteeService *GuaranteeService) GetCorporationGuarantees(request guaranteedto.GetGuaranteesRequest) ([]guaranteedto.GuaranteeResponse, error) {
	err := guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)
	if err != nil {
		return nil, err
	}

	allowedStatus := guaranteeService.mapToFilterStatuses(request.Status)

	guarantees, err := guaranteeService.guaranteeRepository.FindCorporationGuarantees(guaranteeService.db, request.CorporationID, allowedStatus)
	if err != nil {
		return nil, err
	}
	response := make([]guaranteedto.GuaranteeResponse, len(guarantees))

	for i, guarantee := range guarantees {
		response[i], err = guaranteeService.mapGuaranteeToResponse(guarantee)
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}

func (guaranteeService *GuaranteeService) AddGuarantee(request guaranteedto.CreateGuaranteeRequest) (uint, error) {
	err := guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)
	if err != nil {
		return 0, err
	}

	guarantee, err := guaranteeService.guaranteeRepository.FindCorporationGuaranteeByName(guaranteeService.db, request.CorporationID, request.Name)
	if err != nil {
		return 0, err
	}
	if guarantee != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(guaranteeService.constants.Field.Name, guaranteeService.constants.Tag.AlreadyExist)
		return 0, conflictErrors
	}

	guarantee = &entity.Guarantee{
		CorporationID:  request.CorporationID,
		Name:           request.Name,
		Status:         request.Status,
		GuaranteeType:  enum.GuaranteeType(request.GuaranteeType),
		DurationMonths: request.Duration,
		Description:    request.Description,
	}
	err = guaranteeService.db.WithTransaction(func(tx database.Database) error {
		if err := guaranteeService.guaranteeRepository.CreateGuarantee(tx, guarantee); err != nil {
			return err
		}
		for _, terms := range request.GuaranteeTermsRequest {
			if err := guaranteeService.addGuaranteeTerm(terms, guarantee.ID); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return guarantee.ID, nil
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

func (guaranteeService *GuaranteeService) UpdateGuaranteeStatus(request guaranteedto.ChangeStatusRequest) error {
	err := guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)
	if err != nil {
		return err
	}

	guarantee, err := guaranteeService.guaranteeRepository.FindGuaranteeByID(guaranteeService.db, request.GuaranteeID)
	if err != nil {
		return err
	}
	if guarantee == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.Guarantee}
		return notFoundError
	}

	if !enum.GuaranteeStatus(request.Status).IsValid() {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.Invalid)
		return conflictErrors
	}

	if guarantee.Status == enum.GuaranteeStatus(request.Status) {
		var conflictErrors exception.ConflictErrors
		switch guarantee.Status {
		case enum.GuaranteeStatusActive:
			conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.AlreadyActive)
			return conflictErrors
		case enum.GuaranteeStatusArchive:
			conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.AlreadyArchived)
			return conflictErrors
		default:
			conflictErrors.Add(guaranteeService.constants.Field.Guarantee, guaranteeService.constants.Tag.StatusNotChange)
			return conflictErrors
		}
	}

	guarantee.Status = enum.GuaranteeStatus(request.Status)

	if err := guaranteeService.guaranteeRepository.UpdateGuarantee(guaranteeService.db, guarantee); err != nil {
		return err
	}
	return nil
}

func (guaranteeService *GuaranteeService) CreateGuaranteeViolation(request guaranteedto.CreateGuaranteeViolationRequest) (uint, error) {
	err := guaranteeService.corporationService.CheckApplicantAccess(request.CorporationID, request.OperatorID)
	if err != nil {
		return 0, err
	}

	violation, err := guaranteeService.guaranteeRepository.FindPanelGuaranteeViolation(guaranteeService.db, request.PanelID)
	if err != nil {
		return 0, err
	}
	if violation != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(guaranteeService.constants.Field.GuaranteeViolation, guaranteeService.constants.Tag.AlreadyExist)
		return 0, conflictErrors
	}

	violation = &entity.GuaranteeViolation{
		PanelID:      request.PanelID,
		ViolatedByID: request.OperatorID,
		Reason:       request.Reason,
		Details:      request.Details,
	}
	if err := guaranteeService.guaranteeRepository.CreateGuaranteeViolation(guaranteeService.db, violation); err != nil {
		return 0, err
	}

	return violation.ID, nil
}

func (guaranteeService *GuaranteeService) GetCorporationPanelGuaranteeViolation(panelID uint) (guaranteedto.CorporationGuaranteeViolationResponse, error) {
	var response guaranteedto.CorporationGuaranteeViolationResponse

	violation, err := guaranteeService.guaranteeRepository.FindPanelGuaranteeViolation(guaranteeService.db, panelID)
	if err != nil {
		return guaranteedto.CorporationGuaranteeViolationResponse{}, err
	}
	if violation == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.GuaranteeViolation}
		return response, notFoundError
	}

	operator, err := guaranteeService.userService.GetUserCredential(violation.ViolatedByID)
	if err != nil {
		return guaranteedto.CorporationGuaranteeViolationResponse{}, err
	}

	response = guaranteedto.CorporationGuaranteeViolationResponse{
		ViolatedBy: operator,
		Reason:     violation.Reason,
		Details:    violation.Details,
	}

	return response, nil
}

func (guaranteeService *GuaranteeService) GetCustomerPanelGuaranteeViolation(panelID uint) (guaranteedto.CustomerGuaranteeViolationResponse, error) {
	var response guaranteedto.CustomerGuaranteeViolationResponse

	violation, err := guaranteeService.guaranteeRepository.FindPanelGuaranteeViolation(guaranteeService.db, panelID)
	if err != nil {
		return guaranteedto.CustomerGuaranteeViolationResponse{}, err
	}
	if violation == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.GuaranteeViolation}
		return response, notFoundError
	}

	response = guaranteedto.CustomerGuaranteeViolationResponse{
		Reason:  violation.Reason,
		Details: violation.Details,
	}

	return response, nil
}

func (guaranteeService *GuaranteeService) UpdateGuaranteeViolation(request guaranteedto.UpdateGuaranteeViolationRequest) error {
	violation, err := guaranteeService.guaranteeRepository.FindPanelGuaranteeViolation(guaranteeService.db, request.PanelID)
	if err != nil {
		return err
	}
	if violation == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.GuaranteeViolation}
		return notFoundError
	}

	if request.Reason != nil {
		violation.Reason = *request.Reason
	}

	if request.Details != nil {
		violation.Reason = *request.Details
	}

	violation.ViolatedByID = request.OperatorID

	if err := guaranteeService.guaranteeRepository.UpdateGuaranteeViolation(guaranteeService.db, violation); err != nil {
		return err
	}
	return nil
}

func (guaranteeService *GuaranteeService) RemovePanelGuaranteeViolation(panelID uint) error {
	violation, err := guaranteeService.guaranteeRepository.FindPanelGuaranteeViolation(guaranteeService.db, panelID)
	if err != nil {
		return err
	}
	if violation == nil {
		notFoundError := exception.NotFoundError{Item: guaranteeService.constants.Field.GuaranteeViolation}
		return notFoundError
	}

	if err := guaranteeService.guaranteeRepository.DeleteGuaranteeViolation(guaranteeService.db, violation); err != nil {
		return err
	}
	return nil
}
