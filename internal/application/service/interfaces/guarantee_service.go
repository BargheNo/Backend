package service

import guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"

type GuaranteeService interface {
	ValidateGuaranteeOwnerShip(guaranteeID, corporationID uint) error
	GetGuarantee(guaranteeID uint) (guaranteedto.GuaranteeResponse, error)
	GetGuaranteeTypes() []guaranteedto.GuaranteeTypesResponse
	GetGuaranteeStatuses() []guaranteedto.GuaranteeTypesResponse
	GetCorporationGuarantee(request guaranteedto.GetGuaranteeRequest) guaranteedto.GuaranteeResponse
	GetCorporationGuarantees(request guaranteedto.GetGuaranteesRequest) []guaranteedto.GuaranteeResponse
	AddGuarantee(request guaranteedto.CreateGuaranteeRequest) uint
	UpdateGuaranteeStatus(request guaranteedto.ChangeStatusRequest)
}
