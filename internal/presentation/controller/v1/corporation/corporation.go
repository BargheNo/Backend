package corporation

import (
	"mime/multipart"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CorporationCorporationController struct {
	constants          *bootstrap.Constants
	pagination         *bootstrap.Pagination
	corporationService usecase.CorporationService
}

func NewCorporationCorporationController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	corporationService usecase.CorporationService,
) *CorporationCorporationController {
	return &CorporationCorporationController{
		constants:          constants,
		pagination:         pagination,
		corporationService: corporationService,
	}
}

func (corporationController *CorporationCorporationController) GetMyProfile(ctx *gin.Context) {
	type getCorporationParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getCorporationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	corporationRequest := corporationdto.CorporationDetailsRequest{
		UserID:        userID.(uint),
		CorporationID: params.CorporationID,
		Status:        enum.CorpStatusApproved,
	}
	corporationDetails, err := corporationController.corporationService.GetCorporationDetails(corporationRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", corporationDetails)
}

func (corporationController *CorporationCorporationController) GetPublicProfile(ctx *gin.Context) {
	type getCorporationParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
	}
	params := controller.Validated[getCorporationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	corporationRequest := corporationdto.CorporationDetailsRequest{
		CorporationID: params.CorporationID,
		UserID:        userID.(uint),
	}
	corporationDetails, err := corporationController.corporationService.GetCorporationPublicDetails(corporationRequest)
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", corporationDetails)
}

func (corporationController *CorporationCorporationController) AddAddress(ctx *gin.Context) {
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
		CorporationStatus: enum.CorpStatusApproved,
		Addresses:         addresses,
	}

	if err := corporationController.corporationService.AddAddress(addressInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationCorporationController) DeleteAddress(ctx *gin.Context) {
	type deleteAddressParams struct {
		CorporationID uint `uri:"corporationID" validate:"required"`
		AddressID     uint `uri:"addressID" validate:"required"`
	}
	params := controller.Validated[deleteAddressParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	addressInfo := corporationdto.DeleteAddressRequest{
		UserID:            userID.(uint),
		CorporationID:     params.CorporationID,
		CorporationStatus: enum.CorpStatusApproved,
		AddressID:         params.AddressID,
	}
	if err := corporationController.corporationService.DeleteAddress(addressInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationCorporationController) AddContactInformation(ctx *gin.Context) {
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
		CorporationStatus:  enum.CorpStatusApproved,
		ContactInformation: contacts,
	}
	if err := corporationController.corporationService.AddContactInfo(contactInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateContactInfo")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationCorporationController) DeleteContactInformation(ctx *gin.Context) {
	type contactInformationParams struct {
		CorporationID        uint `uri:"corporationID" validate:"required"`
		ContactInformationID uint `uri:"contactID" validate:"required"`
	}
	params := controller.Validated[contactInformationParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	contactInfo := corporationdto.DeleteContactInformationRequest{
		ApplicantID:       userID.(uint),
		ContactID:         params.ContactInformationID,
		CorporationID:     params.CorporationID,
		CorporationStatus: enum.CorpStatusApproved,
	}
	if err := corporationController.corporationService.DeleteContactInfo(contactInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteContactInfo")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationCorporationController) ChangeLogo(ctx *gin.Context) {
	type profileLogoParams struct {
		CorporationID uint                  `uri:"corporationID" validate:"required"`
		Logo          *multipart.FileHeader `form:"logo" validate:"required"`
	}
	params := controller.Validated[profileLogoParams](ctx)
	userID, _ := ctx.Get(corporationController.constants.Context.ID)

	changeLogoRequest := corporationdto.ChangeLogoRequest{
		ApplicantID:   userID.(uint),
		CorporationID: params.CorporationID,
		Logo:          params.Logo,
	}
	if err := corporationController.corporationService.ChangeLogo(changeLogoRequest); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.changeLogo")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationCorporationController) UpdateRegister(ctx *gin.Context) {
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

	if err := corporationController.corporationService.UpdateRegistrationInfoProfile(updateRegisterInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateCorporation")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *CorporationCorporationController) SubmitCertificateFiles(ctx *gin.Context) {
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
	if err := corporationController.corporationService.AddCertificateFilesFromProfile(requestInfo); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addCorporationCertificate")
	controller.Response(ctx, 200, message, nil)
}

// func (corporationController *CorporationCorporationController) CreateCorporationStaff(ctx *gin.Context) {
// 	type addStaffParams struct {
// 		CorporationID uint   `uri:"corporationID" validate:"required"`
// 		Phone         string `json:"phone" validate:"required,e164"`
// 		Role          uint   `json:"role" validate:"role"`
// 	}
// 	params := controller.Validated[addStaffParams](ctx)
// 	userID, _ := ctx.Get(corporationController.constants.Context.ID)

// }
