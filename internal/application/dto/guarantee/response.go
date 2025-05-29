package guaranteedto

type GuaranteeResponse struct {
	ID             uint                    `json:"id"`
	Name           string                  `json:"string"`
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
	ID   uint
	Name string
}
