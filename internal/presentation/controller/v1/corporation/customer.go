package corporation

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerCorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService service.CorporationService
}

func NewCustomerCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService service.CorporationService,
) *CustomerCorporationController {
	return &CustomerCorporationController{
		constants:          constants,
		pagination:         pagination,
		corporationService: corporationService,
	}
}

func (corporationController *CustomerCorporationController) Register(ctx *gin.Context) {
	type signatory struct {
		Name               string `json:"name" validate:"required"`
		NationalCardNumber string `json:"nationalCardNumber" validate:"required"`
		Position           string `json:"position"`
	}
	type registerParams struct {
		Name               string      `json:"name" validate:"required"`
		RegistrationNumber string      `json:"registrationNumber" validate:"required"`
		NationalID         string      `json:"nationalID" validate:"required"`
		IBAN               string      `json:"iban"`
		Signatories        []signatory `json:"signatories" validate:"required"`
	}
	params := controller.Validated[registerParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)
	signatories := make([]corporationdto.Signatory, len(params.Signatories))
	for i, signatory := range params.Signatories {
		signatories[i] = corporationdto.Signatory{
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
	}
	registerInfo := corporationdto.RegisterRequest{
		ApplicantID:        userID.(uint),
		Name:               params.Name,
		NationalID:         params.NationalID,
		RegistrationNumber: params.RegistrationNumber,
		IBAN:               params.IBAN,
		Signatories:        signatories,
	}

	corporationInfo := corporationController.corporationService.Register(registerInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.corporationRegister")
	controller.Response(ctx, 200, message, corporationInfo)
}

func (corporationController *CustomerCorporationController) UpdateRegister(ctx *gin.Context) {
	type signatory struct {
		Name               string `json:"name" validate:"required"`
		NationalCardNumber string `json:"nationalCardNumber" validate:"required"`
		Position           string `json:"position" validate:"required"`
	}

	type registerParams struct {
		CorporationID      uint        `uri:"corporationID" validate:"required"`
		Name               *string     `json:"name"`
		RegistrationNumber *string     `json:"registrationNumber"`
		NationalID         *string     `json:"nationalID"`
		IBAN               *string     `json:"iban"`
		Signatories        []signatory `json:"signatories" validate:"omitempty,dive"`
	}
	params := controller.Validated[registerParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	signatories := make([]corporationdto.Signatory, len(params.Signatories))
	for i, signatory := range params.Signatories {
		signatories[i] = corporationdto.Signatory{
			Name:               signatory.Name,
			NationalCardNumber: signatory.NationalCardNumber,
			Position:           signatory.Position,
		}
	}
	updateRegisterInfo := corporationdto.UpdateRegisterRequest{
		ApplicantID:        userID.(uint),
		CorporationID:      params.CorporationID,
		Name:               params.Name,
		NationalID:         params.NationalID,
		RegistrationNumber: params.RegistrationNumber,
		IBAN:               params.IBAN,
		Signatories:        signatories,
	}

	corporationController.corporationService.UpdateRegister(updateRegisterInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateCorporation")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) GetCorporationPrivateDetails(ctx *gin.Context) {
	type GetCorporationParamsParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[GetCorporationParamsParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)
	corporationRequest := corporationdto.CorporationDetailsRequest{
		UserID:        userID.(uint),
		CorporationID: params.CorporationID,
		Status:        enum.CorpStatusAwaitingApproval,
	}

	corporationDetails := corporationController.corporationService.GetCorporationDetails(corporationRequest)
	controller.Response(ctx, 200, "", corporationDetails)
}

func (corporationController *CustomerCorporationController) AddAddress(ctx *gin.Context) {
	type address struct {
		ProvinceID    uint   `json:"provinceID" validate:"required"`
		CityID        uint   `json:"cityID" validate:"required"`
		StreetAddress string `json:"streetAddress" validate:"required"`
		PostalCode    string `json:"postalCode" validate:"required"`
		HouseNumber   string `json:"houseNumber" validate:"required"`
		Unit          uint   `json:"unit" validate:"required"`
	}
	type addressesInformationParams struct {
		CorporationID uint      `uri:"corporationID" validate:"required"`
		Addresses     []address `json:"addresses" validate:"required"`
	}
	params := controller.Validated[addressesInformationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	addresses := make([]addressdto.CreateAddressRequest, len(params.Addresses))
	for i, address := range params.Addresses {
		addresses[i] = addressdto.CreateAddressRequest{
			ProvinceID:    address.ProvinceID,
			CityID:        address.CityID,
			StreetAddress: address.StreetAddress,
			PostalCode:    address.PostalCode,
			HouseNumber:   address.HouseNumber,
			Unit:          address.Unit,
			OwnerID:       params.CorporationID,
			OwnerType:     corporationController.constants.AddressOwners.Corporation,
		}
	}

	addressInfo := corporationdto.AddCorporationAddressRequest{
		ApplicantID:       userID.(uint),
		CorporationID:     params.CorporationID,
		CorporationStatus: enum.CorpStatusAwaitingApproval,
		Addresses:         addresses,
	}

	corporationController.corporationService.AddAddress(addressInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) DeleteAddress(ctx *gin.Context) {
	type deleteAddressParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		AddressID     uint `uri:"addressID" validate:"required"`
	}
	params := controller.Validated[deleteAddressParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	addressInfo := corporationdto.DeleteAddressRequest{
		UserID:            userID.(uint),
		CorporationID:     params.CorporationID,
		CorporationStatus: enum.CorpStatusAwaitingApproval,
		AddressID:         params.AddressID,
	}
	corporationController.corporationService.DeleteAddress(addressInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) AddContactInformation(ctx *gin.Context) {
	type contactInformation struct {
		ContactTypeID uint   `json:"contactTypeID"`
		ContactValue  string `json:"contactValue"`
	}
	type contactInformationParams struct {
		CorporationID      uint                 `uri:"corporationID" validate:"required"`
		ContactInformation []contactInformation `json:"contactInformation" validate:"required"`
	}
	params := controller.Validated[contactInformationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	contacts := make([]corporationdto.ContactInformation, len(params.ContactInformation))
	for i, contact := range params.ContactInformation {
		contacts[i] = corporationdto.ContactInformation{
			ContactTypeID: contact.ContactTypeID,
			ContactValue:  contact.ContactValue,
		}
	}
	contactInfo := corporationdto.AddContactInformationRequest{
		ApplicantID:        userID.(uint),
		CorporationID:      params.CorporationID,
		CorporationStatus:  enum.CorpStatusAwaitingApproval,
		ContactInformation: contacts,
	}
	corporationController.corporationService.AddContactInfo(contactInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateContactInfo")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) DeleteContactInformation(ctx *gin.Context) {
	type contactInformationParams struct {
		CorporationID        uint `uri:"corporationID" validate:"required"`
		ContactInformationID uint `uri:"contactID" validate:"required"`
	}
	params := controller.Validated[contactInformationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	contactInfo := corporationdto.DeleteContactInformationRequest{
		ApplicantID:       userID.(uint),
		CorporationID:     params.CorporationID,
		CorporationStatus: enum.CorpStatusAwaitingApproval,
	}

	corporationController.corporationService.DeleteContactInfo(contactInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteContactInfo")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) SubmitCertificateFiles(ctx *gin.Context) {
	type certificatesParams struct {
		CorporationID          uint                  `uri:"corporationID" validate:"required"`
		VATTaxpayerCertificate *multipart.FileHeader `form:"vatTaxpayerCertificate"`
		OfficialNewspaperAD    *multipart.FileHeader `form:"officialNewspaperAD"`
	}
	params := controller.Validated[certificatesParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	requestInfo := corporationdto.AddCertificatesRequest{
		CorporationID:          params.CorporationID,
		ApplicantID:            userID.(uint),
		VATTaxpayerCertificate: params.VATTaxpayerCertificate,
		OfficialNewspaperAD:    params.OfficialNewspaperAD,
	}
	corporationController.corporationService.AddCertificateFiles(requestInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addCorporationCertificate")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CustomerCorporationController) GetCorporations(ctx *gin.Context) {
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	pagination := controller.GetPagination(ctx, corporationController.pagination.DefaultPage, corporationController.pagination.DefaultPageSize)
	offset, limit := pagination.GetOffsetLimit()
	corporationRequest := corporationdto.CorporationListRequest{
		UserID: userID.(uint),
		Offset: offset,
		Limit:  limit,
	}

	corporations := corporationController.corporationService.GetCorporations(corporationRequest)
	controller.Response(ctx, 200, "", corporations)
}
