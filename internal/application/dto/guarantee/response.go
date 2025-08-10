package guaranteedto

import userdto "github.com/BargheNo/Backend/internal/application/dto/user"

type GuaranteeResponse struct {
	ID             uint                    `json:"id"`
	Name           string                  `json:"name"`
	Status         string                  `json:"status"`
	GuaranteeType  string                  `json:"guaranteeType"`
	DurationMonths uint                    `json:"durationMonths"`
	Description    string                  `json:"description"`
	Terms          []GuaranteeTermResponse `json:"terms"`
}

type GuaranteeTermResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Limitations string `json:"limitations,omitempty"`
}

type GuaranteeTypesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CorporationGuaranteeViolationResponse struct {
	ViolatedBy userdto.CredentialResponse `json:"operator"`
	Reason     string                     `json:"reason"`
	Details    string                     `json:"details"`
}

type CustomerGuaranteeViolationResponse struct {
	Reason  string `json:"reason"`
	Details string `json:"details"`
}
