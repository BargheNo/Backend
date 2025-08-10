package usecase

import guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"

type GuaranteeService interface {
	ValidateActiveGuaranteeOwnerShip(guaranteeID, corporationID uint) error
	GetGuarantee(guaranteeID uint) (guaranteedto.GuaranteeResponse, error)
	GetGuaranteeTypes() []guaranteedto.GuaranteeTypesResponse
	GetGuaranteeStatuses() []guaranteedto.GuaranteeTypesResponse
	GetCorporationGuarantee(request guaranteedto.GetGuaranteeRequest) (guaranteedto.GuaranteeResponse, error)
	GetCorporationGuarantees(request guaranteedto.GetGuaranteesRequest) ([]guaranteedto.GuaranteeResponse, error)
	AddGuarantee(request guaranteedto.CreateGuaranteeRequest) (uint, error)
	UpdateGuaranteeStatus(request guaranteedto.ChangeStatusRequest) error
	CreateGuaranteeViolation(request guaranteedto.CreateGuaranteeViolationRequest) (uint, error)
	GetCorporationPanelGuaranteeViolation(panelID uint) (guaranteedto.CorporationGuaranteeViolationResponse, error)
	GetCustomerPanelGuaranteeViolation(panelID uint) (guaranteedto.CustomerGuaranteeViolationResponse, error)
	UpdateGuaranteeViolation(request guaranteedto.UpdateGuaranteeViolationRequest) error
	RemovePanelGuaranteeViolation(panelID uint) error
}
