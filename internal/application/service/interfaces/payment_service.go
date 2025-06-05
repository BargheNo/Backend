package service

import paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"

type PaymentService interface {
	GetPaymentTerms(payTermID uint) (paymentdto.PaymentTermsResponse, error)
	CreatePaymentTerms(paymentTermsRequest paymentdto.PaymentTermsRequest) uint
	UpdatePaymentTerms(updatePaymentRequest paymentdto.UpdatePaymentTermsRequest) error
}
