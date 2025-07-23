package usecase

import (
	guaranteedto "github.com/BargheNo/Backend/internal/application/dto/guarantee"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
)

type InstallationService interface {
	AddPanel(panelInfo installationdto.AddPanelRequest) error
	ChangeInstallationRequestStatus(request installationdto.ChangeRequestStatusRequest) error
	ClearPanelGuaranteeViolation(violationInfo installationdto.GetCorporationGuaranteeViolationRequest) error
	CompleteInstallationRequest(request installationdto.CompleteBidInstallationRequest) error
	CreateInstallationRequest(request installationdto.NewInstallationRequest) error
	GetAnonymousInstallationRequest(request installationdto.CorporationPanelRequest) (installationdto.AnonymousRequestsResponse, error)
	GetAnonymousInstallationRequests(request installationdto.CorporationPanelListRequest) ([]installationdto.AnonymousRequestsResponse, error)
	GetBuildingTypes() []installationdto.EnumStatusResponse
	GetCorporationPanel(request installationdto.CorporationPanelRequest) (installationdto.CorporationPanelResponse, error)
	GetCorporationPanelGuaranteeViolation(violationInfo installationdto.GetCorporationGuaranteeViolationRequest) (guaranteedto.CorporationGuaranteeViolationResponse, error)
	GetCorporationPanels(listInfo installationdto.CorporationPanelListRequest) ([]installationdto.CorporationPanelListResponse, error)
	GetCustomerPanel(panelInfo installationdto.GetOwnerRequest) (installationdto.CustomerPanelResponse, error)
	GetCustomerPanelGuaranteeViolation(violationInfo installationdto.GetCustomerGuaranteeViolationRequest) (guaranteedto.CustomerGuaranteeViolationResponse, error)
	GetCustomerPanels(listInfo installationdto.CustomerPanelListRequest) ([]installationdto.CustomerPanelListResponse, error)
	GetOwnerInstallationRequest(request installationdto.GetOwnerRequest) (installationdto.AnonymousRequestsResponse, error)
	GetOwnerInstallationRequests(request installationdto.CustomerRequestsListRequest) ([]installationdto.AnonymousRequestsResponse, error)
	DeleteInstallationRequest(requestID uint) error
	GetPanelByAdmin(panelID uint) (installationdto.AdminPanelResponse, error)
	GetPanelsByAdmin(listInfo installationdto.AdminInstallationListRequest) ([]installationdto.AdminPanelResponse, error)
	GetPublicInstallationRequest(requestID uint) (installationdto.PublicRequestDetailsResponse, error)
	GetInstallationRequestsByAdmin(request installationdto.AdminInstallationListRequest) ([]installationdto.PublicRequestDetailsResponse, error)
	UpdatePanel(request installationdto.UpdatePanelRequest) error
	DeletePanel(panelID uint) error
	GetRequestStatuses() []installationdto.EnumStatusResponse
	GetPanelStatuses() []installationdto.EnumStatusResponse
	UpdateInstallationRequestByAdmin(newRequest installationdto.UpdateInstallationRequest) error
	UpdatePanelGuaranteeViolation(violationInfo installationdto.UpdateGuaranteeViolationRequest) error
	GetPanelStatus() []installationdto.EnumStatusResponse
	ValidatePanelGuarantee(panelID uint) error
	ValidatePanelOwnership(panelID uint, userID uint) (installationdto.AdminPanelResponse, error)
	ValidateRequestOwnership(requestID uint, ownerID uint) (installationdto.PublicRequestDetailsResponse, error)
	ViolatePanelGuaranteeStatus(request installationdto.CreateViolatePanelGuaranteeRequest) (uint, error)
}
