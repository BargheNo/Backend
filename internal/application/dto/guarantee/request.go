package guaranteedto

import "github.com/BargheNo/Backend/internal/domain/enum"

type GuaranteeTermsRequest struct {
	Title       string
	Description string
	Limitations string
}
type CreateGuaranteeRequest struct {
	CorporationID         uint
	OperatorID            uint
	GuaranteeType         uint
	Name                  string
	Status                enum.GuaranteeStatus
	Duration              uint
	Description           string
	GuaranteeTermsRequest []GuaranteeTermsRequest
}

type GetGuaranteesRequest struct {
	CorporationID uint
	OperatorID    uint
	Status        uint
}

type GetGuaranteeRequest struct {
	CorporationID uint
	OperatorID    uint
	GuaranteeID   uint
}

type ChangeStatusRequest struct {
	CorporationID uint
	OperatorID    uint
	GuaranteeID   uint
	Status        uint
}

type CreateGuaranteeViolationRequest struct {
	PanelID       uint
	CorporationID uint
	OperatorID    uint
	Reason        string
	Details       string
}

type UpdateGuaranteeViolationRequest struct {
	PanelID    uint
	OperatorID uint
	Reason     *string
	Details    *string
}
