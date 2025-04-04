package service

import (
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	"github.com/BargheNo/Backend/internal/domain/entity"
)

type InstallationService interface {
	GetInstallationRequestModel(requestID uint) *entity.InstallationRequest
	CreateInstallationRequest(requestInfo installationdto.NewInstallationRequest)
	GetOwnerInstallationRequests(listInfo installationdto.InstallationListRequest) []installationdto.OwnerRequestsResponse
	GetInstallationRequest(requestID uint) installationdto.RequestDetailsResponse
	GetOwnerInstallationRequest(requestInfo installationdto.GetOwnerRequest) installationdto.OwnerRequestsResponse
	GetInstallationRequests(listInfo installationdto.InstallationListRequest) []installationdto.RequestDetailsResponse
	AddPanel(panelInfo installationdto.AddPanelRequest)
}
