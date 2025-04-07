package bootstrap

import "fmt"

type Constants struct {
	Context       Context
	LogLevel      LogLevel
	RedisKey      RedisKey
	S3BucketPath  BucketPath
	Field         ErrorField
	Tag           ErrorTag
	SMSTemplates  SMSTemplates
	JWTKeysPath   JWTKeysPath
	Metrics       Metrics
	AddressOwners AddressOwners
}

type Context struct {
	Translator                   string
	IsLoadedValidationTranslator string
	ID                           string
}

type LogLevel struct {
	Debug string
	Info  string
	Warn  string
	Error string
	Fatal string
}

type RedisKey struct {
}

type BucketPath struct {
}

type ErrorField struct {
	User                string
	Phone               string
	Password            string
	OTP                 string
	Corporation         string
	NationalID          string
	RegistrationNumber  string
	IBAN                string
	InstallationRequest string
	Bid                 string
	Address             string
	Name                string
	Province            string
	City                string
	Page                string
	ContactType         string
	PanelName           string
	Panel               string
	MaintenanceRequest  string
}

type ErrorTag struct {
	AlreadyRegistered      string
	MinimumLength          string
	ContainsLowercase      string
	ContainsUppercase      string
	ContainsNumber         string
	ContainsSpecialChar    string
	Expired                string
	Invalid                string
	NotRegistered          string
	NotVerified            string
	InvalidAuthCredentials string
	ExpiredAuthToken       string
	InvalidAuthToken       string
	Unauthorized           string
	AwaitingApproval       string
	Rejected               string
	NotExist               string
	AlreadyExist           string
	ForbiddenStatus        string
	Pending                string
}

type SMSTemplates struct {
	OTP string
}

type JWTKeysPath struct {
	PublicKey  string
	PrivateKey string
}

type Metrics struct {
	HTTPRequestsTotal   Options
	HTTPRequestDuration Options
}

type Options struct {
	Name string
	Help string
}

type AddressOwners struct {
	User                string
	Corporation         string
	InstallationRequest string
	Panel               string
	MaintenanceRequest  string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator:                   "translator",
			IsLoadedValidationTranslator: "isLoadedValidationTranslator",
			ID:                           "ID",
		},
		LogLevel: LogLevel{
			Debug: "debug",
			Info:  "info",
			Warn:  "warn",
			Error: "error",
			Fatal: "fatal",
		},
		Field: ErrorField{
			User:                "user",
			Phone:               "phone",
			Password:            "password",
			OTP:                 "otp",
			Corporation:         "corporation",
			NationalID:          "nationalID",
			RegistrationNumber:  "registrationNumber",
			IBAN:                "iban",
			InstallationRequest: "installationRequest",
			Bid:                 "bid",
			Address:             "address",
			Name:                "name",
			Province:            "province",
			City:                "city",
			Page:                "page",
			ContactType:         "contactType",
			PanelName:           "panelName",
			Panel:               "panel",
			MaintenanceRequest:  "maintenanceRequest",
		},
		Tag: ErrorTag{
			AlreadyRegistered:      "alreadyRegistered",
			MinimumLength:          "minimumLength",
			ContainsLowercase:      "containsLowercase",
			ContainsUppercase:      "containsUppercase",
			ContainsNumber:         "containsNumber",
			ContainsSpecialChar:    "containsSpecialChar",
			Expired:                "Expired",
			Invalid:                "invalid",
			NotRegistered:          "notRegistered",
			NotVerified:            "notVerified",
			InvalidAuthCredentials: "invalidAuthCredentials",
			ExpiredAuthToken:       "expiredAuthToken",
			InvalidAuthToken:       "invalidAuthToken",
			Unauthorized:           "unauthorized",
			AwaitingApproval:       "awaitingApproval",
			Rejected:               "rejected",
			NotExist:               "notExist",
			AlreadyExist:           "alreadyExist",
			ForbiddenStatus:        "forbiddenStatus",
			Pending:                "pending",
		},
		SMSTemplates: SMSTemplates{
			OTP: "sendOTPTemplate",
		},
		JWTKeysPath: JWTKeysPath{
			PublicKey:  "./internal/application/adapter/jwt/publicKey.pem",
			PrivateKey: "./internal/application/adapter/jwt/privateKey.pem",
		},
		Metrics: Metrics{
			HTTPRequestsTotal: Options{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			HTTPRequestDuration: Options{
				Name: "http_request_duration_seconds",
				Help: "HTTP request duration in seconds",
			},
		},
		AddressOwners: AddressOwners{
			User:                "users",
			Corporation:         "corporations",
			InstallationRequest: "installation_requests",
			Panel:               "panels",
		},
	}
}

func (r *RedisKey) GenerateOTPKey(value string) string {
	return fmt.Sprintf("otp:%s", value)
}

func (path *BucketPath) GetVATTaxpayerCertificatePath(corporationID uint, certificateFilename string) string {
	return fmt.Sprintf("corporation/%d/taxpayer/%s", corporationID, certificateFilename)
}

func (path *BucketPath) GetOfficialNewspaperADPath(corporationID uint, certificateFilename string) string {
	return fmt.Sprintf("corporation/%d/newspaper-ad/%s", corporationID, certificateFilename)
}
