package paymentdto

type PaymentTermsResponse struct {
	ID              uint                     `json:"id"`
	PaymentMethod   string                   `json:"paymentMethod"`
	InstallmentPlan *InstallmentPlanResponse `json:"installmentPlan,omitempty"`
}

type InstallmentPlanResponse struct {
	ID                uint   `json:"id"`
	NumberOfMonths    uint   `json:"numberOfMonths"`
	DownPaymentAmount uint   `json:"downPaymentAmount"`
	MonthlyAmount     uint   `json:"monthlyAmount"`
	Notes             string `json:"notes,omitempty"`
	// DownPaymentDate   string `json:"downPaymentDate"`
	// DueDay            uint   `json:"dueDay"`
}

type PaymentMethodResponse struct {
	ID     uint   `json:"id"`
	Method string `json:"method"`
}
