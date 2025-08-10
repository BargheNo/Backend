package paymentdto

type InstallmentPlanRequest struct {
	PaymentTermsID    uint
	NumberOfMonths    uint
	DownPaymentAmount uint
	MonthlyAmount     uint
	Notes             string
}

type PaymentTermsRequest struct {
	PaymentMethod   uint
	InstallmentPlan *InstallmentPlanRequest
}

type UpdateInstallmentPlanRequest struct {
	PaymentTermsID    uint
	NumberOfMonths    *uint
	DownPaymentAmount *uint
	MonthlyAmount     *uint
	Notes             *string
}

type UpdatePaymentTermsRequest struct {
	ID              uint
	PaymentMethod   *uint
	InstallmentPlan *UpdateInstallmentPlanRequest
}
