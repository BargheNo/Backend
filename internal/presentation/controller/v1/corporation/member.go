package corporation

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type MemberCorporationController struct {
	constants          *bootstrap.Constants
	corporationService service.CorporationService
	BidService         service.BidService
}

func NewMemberCorporationController(
	constants *bootstrap.Constants,
	corporationService service.CorporationService,
	BidService service.BidService,
) *MemberCorporationController {
	return &MemberCorporationController{
		constants:          constants,
		corporationService: corporationService,
		BidService:         BidService,
	}
}

func (corporationController *MemberCorporationController) UpdateContactInfo(ctx *gin.Context) {
	type addContactInfoParams struct {
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		Eitaa     string `json:"eitaa"`
		Bale      string `json:"bale"`
		Website   string `json:"website"`
		WhatsApp  string `json:"whatsApp"`
		Instagram string `json:"instagram"`
		Telegram  string `json:"telegram"`
		Linkedin  string `json:"linkedin"`
	}

	params := controller.Validated[addContactInfoParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	contactInfo := corporationdto.ContactInfoRequest{
		Phone:     params.Phone,
		Email:     params.Email,
		Eitaa:     params.Eitaa,
		Bale:      params.Bale,
		Website:   params.Website,
		WhatsApp:  params.WhatsApp,
		Instagram: params.Instagram,
		Telegram:  params.Telegram,
		Linkedin:  params.Linkedin,
	}

	corporationController.corporationService.UpdateContactInfo(corporationID.(uint), contactInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.UpdateContactInfo")
	controller.Response(ctx, 200, message, nil)

}

func (corporationController *MemberCorporationController) AddAddress(ctx *gin.Context) {
	type addAddressParams struct {
		Province       string `json:"province" validate:"required"`
		City           string `json:"city" validate:"required"`
		StreetAddress  string `json:"streetAddress" validate:"required"`
		PostalCode     string `json:"postalCode" validate:"required"`
		BuildingNumber string `json:"buildingNumber" validate:"required"`
		Unit           uint   `json:"unitNumber"`
	}
	params := controller.Validated[addAddressParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	address := corporationdto.AddressRequest{
		Province:       params.Province,
		City:           params.City,
		StreetAddress:  params.StreetAddress,
		PostalCode:     params.PostalCode,
		BuildingNumber: params.BuildingNumber,
		Unit:           params.Unit,
	}
	corporationController.corporationService.AddAddress(corporationID.(uint), address)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.AddAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) EditAddress(ctx *gin.Context) {
	type editAddressParams struct {
		AddressID      uint   `json:"addressId" validate:"required"`
		Province       string `json:"province" validate:"required"`
		City           string `json:"city" validate:"required"`
		StreetAddress  string `json:"streetAddress" validate:"required"`
		PostalCode     string `json:"postalCode" validate:"required"`
		BuildingNumber string `json:"buildingNumber" validate:"required"`
		Unit           uint   `json:"unitNumber"`
	}
	params := controller.Validated[editAddressParams](ctx, &corporationController.constants.Context)

	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	address := corporationdto.AddressRequest{
		Province:       params.Province,
		City:           params.City,
		StreetAddress:  params.StreetAddress,
		PostalCode:     params.PostalCode,
		BuildingNumber: params.BuildingNumber,
		Unit:           params.Unit,
	}
	corporationController.corporationService.EditAddress(corporationID.(uint), params.AddressID, address)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.EditAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) DeleteAddress(ctx *gin.Context) {
	type deleteAddressParams struct {
		AddressID uint `json:"addressId" validate:"required"`
	}
	params := controller.Validated[deleteAddressParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	corporationController.corporationService.DeleteAddress(corporationID.(uint), params.AddressID)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.DeleteAddress")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) GetInstallationRequests(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, &corporationController.constants.Context)
	sortparams := controller.GetSort(ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	installation_requests := corporationController.BidService.GetInstallationRequests(corporationID.(uint), pagination.Page, pagination.PageSize, sortparams.SortBy, sortparams.Ascending)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetInstallationRequests")
	controller.Response(ctx, 200, message, installation_requests)
}

func (corporationController *MemberCorporationController) SetBid(ctx *gin.Context) {
	type setBidParams struct {
		InstallationRequestID uint    `json:"installationRequestId" validate:"required"`
		MinCost               float64 `json:"minCost"`
		MaxCost               float64 `json:"maxCost"`
		MinDeadline           time.Time  `json:"minDeadline"`
		MaxDeadline           time.Time  `json:"maxDeadline"`
		Description           string  `json:"description"`
		InstallationTime      string  `json:"installationTime" `
	}
	params := controller.Validated[setBidParams](ctx, &corporationController.constants.Context)

	bidInfo := corporationdto.SetBidRequest{
		InstallationRequestID: params.InstallationRequestID,
		MinCost:               params.MinCost,
		MaxCost:               params.MaxCost,
		MinDeadline:           params.MinDeadline,
		MaxDeadline:           params.MaxDeadline,
		Description:           params.Description,
		InstallationTime:      params.InstallationTime,
	}

	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo.CorporationID = corporationID.(uint)
	corporationController.BidService.SetBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.SetBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) CancelBid(ctx *gin.Context) {
	type cancelBidParams struct {
		BidID                 uint `json:"bidId" validate:"required"`
		InstallationRequestID uint `json:"installationRequestId" validate:"required"`
	}
	params := controller.Validated[cancelBidParams](ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bidInfo := corporationdto.CancelBidRequest{
		BidID:                 params.BidID,
		InstallationRequestID: params.InstallationRequestID,
		CorporationID:         corporationID.(uint),
	}
	corporationController.BidService.CancelBid(bidInfo)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.CancelBid")
	controller.Response(ctx, 200, message, nil)
}

func (corporationController *MemberCorporationController) GetBids(ctx *gin.Context) {
	pagination := controller.GetPagination(ctx, &corporationController.constants.Context)
	sortparams := controller.GetSort(ctx, &corporationController.constants.Context)
	corporationID, _ := ctx.Get(corporationController.constants.Context.ID)
	bids := corporationController.BidService.GetBids(corporationID.(uint), pagination.Page, pagination.PageSize, sortparams.SortBy, sortparams.Ascending)

	trans := controller.GetTranslator(ctx, corporationController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.GetBids")
	controller.Response(ctx, 200, message, bids)
}
