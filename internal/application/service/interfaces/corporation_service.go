package service

import corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"

type CorporationService interface {
	Register(registerInfo corporationdto.RegisterRequest) error
}
