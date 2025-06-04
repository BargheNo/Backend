package service

import guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"

type GuaranteeService interface {
	ValidateActiveGuaranteeOwnerShip(guaranteeID, corporationID uint) error
	GetGuarantee(guaranteeID uint) (guaranteedto.GuaranteeResponse, error)
	GetGuaranteeTypes() []guaranteedto.GuaranteeTypesResponse
	GetGuaranteeStatuses() []guaranteedto.GuaranteeTypesResponse
	GetCorporationGuarantee(request guaranteedto.GetGuaranteeRequest) guaranteedto.GuaranteeResponse
	GetCorporationGuarantees(request guaranteedto.GetGuaranteesRequest) []guaranteedto.GuaranteeResponse
	AddGuarantee(request guaranteedto.CreateGuaranteeRequest) uint
	UpdateGuaranteeStatus(request guaranteedto.ChangeStatusRequest)
	CreateGuaranteeViolation(request guaranteedto.CreateGuaranteeViolationRequest) uint
	GetCorporationPanelGuaranteeViolation(panelID uint) (guaranteedto.CorporationGuaranteeViolationResponse, error)
	GetCustomerPanelGuaranteeViolation(panelID uint) (guaranteedto.CustomerGuaranteeViolationResponse, error)
	UpdateGuaranteeViolation(request guaranteedto.UpdateGuaranteeViolationRequest)
	RemovePanelGuaranteeViolation(panelID uint)
}
