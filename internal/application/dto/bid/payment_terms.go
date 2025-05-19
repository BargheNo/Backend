package biddto

type InstallmentPlanRequest struct {
	NumberOfMonths    uint   `json:"numberOfMonths" validate:"required"`
	DownPaymentAmount uint   `json:"downPaymentAmount" validate:"required"`
	MonthlyAmount     uint   `json:"monthlyAmount" validate:"required"`
	DownPaymentDate   string `json:"downPaymentDate" validate:"required"`
	DueDay            uint   `json:"dueDay" validate:"required,min=1,max=31"`
	Notes             string `json:"notes,omitempty"`
}

type PaymentTermsRequest struct {
	PaymentMethod   string                `json:"paymentMethod" validate:"required,oneof=cash installment"`
	InstallmentPlan *InstallmentPlanRequest `json:"installmentPlan,omitempty" validate:"required_if=PaymentMethod installment"`
}

type InstallmentPlanResponse struct {
	NumberOfMonths    uint   `json:"numberOfMonths"`
	DownPaymentAmount uint   `json:"downPaymentAmount"`
	MonthlyAmount     uint   `json:"monthlyAmount"`
	DownPaymentDate   string `json:"downPaymentDate"`
	DueDay            uint   `json:"dueDay"`
	Notes             string `json:"notes,omitempty"`
}

type PaymentTermsResponse struct {
	PaymentMethod   string                  `json:"paymentMethod"`
	InstallmentPlan *InstallmentPlanResponse `json:"installmentPlan,omitempty"`
} 