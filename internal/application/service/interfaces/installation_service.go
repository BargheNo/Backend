package service

import installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"

type InstallationService interface {
	CreateInstallationRequest(requestInfo installationdto.NewInstallationRequest)
	GetOwnerInstallationRequests(listInfo installationdto.ListOwnerRequestsRequest) []installationdto.ListOwnerRequestsResponse
}
