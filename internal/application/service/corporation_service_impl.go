package serviceimpl

import (
	"regexp"

	"github.com/BargheNo/Backend/bootstrap"
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enums"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
)

type CorporationService struct {
	constants             *bootstrap.Constants
	JWTService            service.JWTService
	db                    database.Database
	CorporationRepository repository.CorporationRepository
	CINService            service.CINService
}

func NewCorporationService(
	constants *bootstrap.Constants,
	jwtService service.JWTService,
	db database.Database,
	corporationRepository repository.CorporationRepository,
	cinService service.CINService,
) *CorporationService {
	return &CorporationService{
		constants:             constants,
		JWTService:            jwtService,
		db:                    db,
		CorporationRepository: corporationRepository,
		CINService:            cinService,
	}
}

func (corporationService *CorporationService) GetCorporationByID(corporationID uint) (*entity.Corporation, bool) {
	corporation, exist := corporationService.CorporationRepository.FindCorporationByID(corporationService.db, corporationID)
	if !exist {
		return nil, false
	}
	if corporation.Status != enums.Approved {
		return nil, false
	}
	return corporation, true
}

func (corporationService *CorporationService) validatePasswordTests(errors *[]string, test string, password string, tag string) {
	matched, _ := regexp.MatchString(test, password)
	if !matched {
		*errors = append(*errors, tag)
	}
}

func (corporationService *CorporationService) passwordValidation(password string) error {
	var errors exception.ValidationErrors
	var errorTags []string

	corporationService.validatePasswordTests(&errorTags, ".{8,}", password, corporationService.constants.Tag.MinimumLength)
	corporationService.validatePasswordTests(&errorTags, "[a-z]", password, corporationService.constants.Tag.ContainsLowercase)
	corporationService.validatePasswordTests(&errorTags, "[A-Z]", password, corporationService.constants.Tag.ContainsUppercase)
	corporationService.validatePasswordTests(&errorTags, "[0-9]", password, corporationService.constants.Tag.ContainsNumber)
	corporationService.validatePasswordTests(&errorTags, "[^\\d\\w]", password, corporationService.constants.Tag.ContainsSpecialChar)

	for _, tag := range errorTags {
		errors.Add(corporationService.constants.Field.Password, tag)
	}
	if len(errorTags) > 0 {
		return errors
	}

	return nil
}

func (corporationService *CorporationService) Register(registerInfo corporationdto.RegisterRequest) {
	var conflictErrors exception.ConflictErrors
	corporation, exist := corporationService.CorporationRepository.FindCorporationByCIN(corporationService.db, registerInfo.CIN)
	if exist && corporation.Status != enums.Rejected {
		conflictErrors.Add(corporationService.constants.Field.CIN, corporationService.constants.Tag.AlreadyRegistered)
		panic(conflictErrors)
	}
	_, err := corporationService.CINService.ValidateCIN(registerInfo.CIN)
	if err != nil {
		panic(err)
	}

	err = corporationService.passwordValidation(registerInfo.Password)
	if err != nil {
		panic(err)
	}

	hashedPassword, err := hashPassword(registerInfo.Password)
	if err != nil {
		panic(err)
	}

	corporation = &entity.Corporation{
		Name:     registerInfo.Name,
		CIN:      registerInfo.CIN,
		Password: hashedPassword,
		Status:   enums.AwaitingApproval,
	}

	err = corporationService.CorporationRepository.CreateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) Login(loginInfo corporationdto.LoginRequest) corporationdto.CorporationInfoResponse {
	var notFoundError exception.NotFoundError
	corporation, exist := corporationService.CorporationRepository.FindCorporationByCIN(corporationService.db, loginInfo.CIN)
	if !exist {
		notFoundError = exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}
	if corporation.Status != enums.Approved {
		notFoundError = exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}

	err := bcrypt.CompareHashAndPassword([]byte(corporation.Password), []byte(loginInfo.Password))
	if err != nil {
		authError := exception.NewInvalidCredentialsError("cin and password not match", nil)
		panic(authError)
	}

	accessToken, refreshToken := corporationService.JWTService.GenerateToken(corporation.ID)
	return corporationdto.CorporationInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         corporation.Name,
	}
}

func (corporationService *CorporationService) ChangePassword(changePasswordRequest corporationdto.ChangePasswordRequest) {
	corporation, exist := corporationService.GetCorporationByID(changePasswordRequest.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}

	err := corporationService.passwordValidation(changePasswordRequest.NewPassword)
	if err != nil {
		panic(err)
	}

	hashedPassword, err := hashPassword(changePasswordRequest.NewPassword)
	if err != nil {
		panic(err)
	}

	corporation.Password = hashedPassword
	err = corporationService.CorporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) UpdateContactInfo(corporationID uint, contactInfo corporationdto.ContactInfoRequest) {
	corporation, exist := corporationService.GetCorporationByID(corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}

	newContactInformation := &entity.ContactInformation{
		Phone:     contactInfo.Phone,
		Email:     contactInfo.Email,
		Eitaa:     contactInfo.Eitaa,
		Bale:      contactInfo.Bale,
		Website:   contactInfo.Website,
		WhatsApp:  contactInfo.WhatsApp,
		Instagram: contactInfo.Instagram,
		Telegram:  contactInfo.Telegram,
		Linkedin:  contactInfo.Linkedin,
	}
	corporation.ContactInformation = *newContactInformation
	err := corporationService.CorporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) AddAddress(address corporationdto.AddressRequest) {
	corporation, exist := corporationService.GetCorporationByID(address.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}

	newAddress := &entity.Address{
		Province:       address.Province,
		City:           address.City,
		StreetAddress:  address.StreetAddress,
		PostalCode:     address.PostalCode,
		BuildingNumber: address.BuildingNumber,
		Unit:           address.Unit,
	}

	corporation.Addresses = append(corporation.Addresses, *newAddress)

	err := corporationService.CorporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) EditAddress(addressID uint, address corporationdto.AddressRequest) {
	corporation, exist := corporationService.GetCorporationByID(address.CorporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}

	newAddress := &entity.Address{
		Province:       address.Province,
		City:           address.City,
		StreetAddress:  address.StreetAddress,
		PostalCode:     address.PostalCode,
		BuildingNumber: address.BuildingNumber,
		Unit:           address.Unit,
	}

	for i, addr := range corporation.Addresses {
		if addr.ID == addressID {
			corporation.Addresses[i] = *newAddress
			break
		}
	}

	err := corporationService.CorporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}

func (corporationService *CorporationService) DeleteAddress(corporationID uint, addressID uint) {
	corporation, exist := corporationService.GetCorporationByID(corporationID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: corporationService.constants.Field.Corporation}
		panic(notFoundError)
	}

	for i, addr := range corporation.Addresses {
		if addr.ID == addressID {
			corporation.Addresses = append(corporation.Addresses[:i], corporation.Addresses[i+1:]...)
			break
		}
	}
	err := corporationService.CorporationRepository.UpdateCorporation(corporationService.db, corporation)
	if err != nil {
		panic(err)
	}
}
