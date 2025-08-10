package service

import (
	"github.com/BargheNo/Backend/bootstrap"
	paymentdto "github.com/BargheNo/Backend/internal/application/dto/payment"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type PaymentService struct {
	constants         *bootstrap.Constants
	paymentRepository postgres.PaymentRepository
	db                database.Database
}

func NewPaymentService(
	constants *bootstrap.Constants,
	paymentRepository postgres.PaymentRepository,
	db database.Database,
) *PaymentService {
	return &PaymentService{
		constants:         constants,
		paymentRepository: paymentRepository,
		db:                db,
	}
}

func (paymentService *PaymentService) GetPaymentMethods() []paymentdto.PaymentMethodResponse {
	methods := enum.GetAllPaymentMethods()
	response := make([]paymentdto.PaymentMethodResponse, len(methods))
	for i, method := range methods {
		response[i] = paymentdto.PaymentMethodResponse{
			ID:     uint(method),
			Method: method.String(),
		}
	}
	return response
}

func (paymentService *PaymentService) GetPaymentTerms(payTermID uint) (paymentdto.PaymentTermsResponse, error) {
	paymentTerms, err := paymentService.paymentRepository.FindPaymentTerms(paymentService.db, payTermID)
	if err != nil {
		return paymentdto.PaymentTermsResponse{}, err
	}
	if paymentTerms == nil {
		notFoundError := exception.NotFoundError{Item: paymentService.constants.Field.PaymentTerm}
		return paymentdto.PaymentTermsResponse{}, notFoundError
	}

	response := paymentdto.PaymentTermsResponse{
		ID:            paymentTerms.ID,
		PaymentMethod: paymentTerms.PaymentMethod.String(),
	}

	if paymentTerms.PaymentMethod == enum.PaymentMethodInstallment {
		installmentPlan, err := paymentService.paymentRepository.FindPaymentTermInstallmentPlan(paymentService.db, payTermID)
		if err != nil {
			return paymentdto.PaymentTermsResponse{}, err
		}
		if installmentPlan == nil {
			notFoundError := exception.NotFoundError{Item: paymentService.constants.Field.PaymentTerm}
			return response, notFoundError
		}
		response.InstallmentPlan = &paymentdto.InstallmentPlanResponse{
			NumberOfMonths:    installmentPlan.NumberOfMonths,
			DownPaymentAmount: installmentPlan.DownPaymentAmount,
			MonthlyAmount:     installmentPlan.MonthlyAmount,
			Notes:             installmentPlan.Notes,
		}
	}
	return response, nil
}

func (paymentService *PaymentService) CreatePaymentTerms(paymentTermsRequest paymentdto.PaymentTermsRequest) (uint, error) {
	terms := &entity.PaymentTerm{
		PaymentMethod: enum.PaymentMethod(paymentTermsRequest.PaymentMethod),
	}
	if err := paymentService.paymentRepository.CreatePaymentTerms(paymentService.db, terms); err != nil {
		return 0, err
	}
	if paymentTermsRequest.InstallmentPlan != nil {
		paymentTermsRequest.InstallmentPlan.PaymentTermsID = terms.ID
		if err := paymentService.createInstallmentPlan(*paymentTermsRequest.InstallmentPlan); err != nil {
			return 0, err
		}
	}
	return terms.ID, nil
}

func (paymentService *PaymentService) createInstallmentPlan(installmentPlan paymentdto.InstallmentPlanRequest) error {
	plan := &entity.InstallmentPlan{
		PaymentTermsID:    installmentPlan.PaymentTermsID,
		NumberOfMonths:    installmentPlan.NumberOfMonths,
		DownPaymentAmount: installmentPlan.DownPaymentAmount,
		MonthlyAmount:     installmentPlan.MonthlyAmount,
		Notes:             installmentPlan.Notes,
	}
	if err := paymentService.paymentRepository.CreateInstallmentPlan(paymentService.db, plan); err != nil {
		return err
	}
	return nil
}

func (paymentService *PaymentService) UpdatePaymentTerms(updatePaymentRequest paymentdto.UpdatePaymentTermsRequest) error {
	terms, err := paymentService.paymentRepository.FindPaymentTerms(paymentService.db, updatePaymentRequest.ID)
	if err != nil {
		return err
	}
	if terms == nil {
		notFoundError := exception.NotFoundError{Item: paymentService.constants.Field.PaymentTerm}
		return notFoundError
	}

	if updatePaymentRequest.PaymentMethod != nil {
		terms.PaymentMethod = enum.PaymentMethod(*updatePaymentRequest.PaymentMethod)
	}

	if updatePaymentRequest.InstallmentPlan != nil {
		updatePaymentRequest.InstallmentPlan.PaymentTermsID = terms.ID
		if err := paymentService.updateInstallmentPlan(*updatePaymentRequest.InstallmentPlan); err != nil {
			return err
		}
	}
	return nil
}

func (paymentService *PaymentService) updateInstallmentPlan(updateInstallmentPlan paymentdto.UpdateInstallmentPlanRequest) error {
	plan, err := paymentService.paymentRepository.FindPaymentTermInstallmentPlan(paymentService.db, updateInstallmentPlan.PaymentTermsID)
	if err != nil {
		return err
	}
	if plan == nil {
		notFoundError := exception.NotFoundError{Item: paymentService.constants.Field.PaymentTerm}
		return notFoundError
	}
	if updateInstallmentPlan.NumberOfMonths != nil {
		plan.NumberOfMonths = *updateInstallmentPlan.NumberOfMonths
	}

	if updateInstallmentPlan.DownPaymentAmount != nil {
		plan.DownPaymentAmount = *updateInstallmentPlan.DownPaymentAmount
	}

	if updateInstallmentPlan.MonthlyAmount != nil {
		plan.MonthlyAmount = *updateInstallmentPlan.MonthlyAmount
	}

	if updateInstallmentPlan.Notes != nil {
		plan.Notes = *updateInstallmentPlan.Notes
	}

	if err := paymentService.paymentRepository.UpdateInstallmentPlan(paymentService.db, plan); err != nil {
		return err
	}
	return nil
}
