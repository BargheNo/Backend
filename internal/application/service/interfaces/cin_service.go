package service

type CINService interface {
	ValidateCIN(cin string) (bool, error)
}