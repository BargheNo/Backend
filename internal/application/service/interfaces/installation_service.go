package service

import installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"

type InstallationService interface {
	InstallationRequest(requestInfo installationdto.NewInstallationRequest)
}
