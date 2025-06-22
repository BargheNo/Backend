package bootstrap

import "fmt"

type Constants struct {
	Context             Context
	LogLevel            LogLevel
	RedisKey            RedisKey
	S3BucketPath        BucketPath
	Field               ErrorField
	Tag                 ErrorTag
	SMSTemplates        SMSTemplates
	EmailTemplates      EmailTemplates
	JWTKeysPath         JWTKeysPath
	Metrics             Metrics
	AddressOwners       AddressOwners
	TicketOwners        TicketOwners
	TicketCommentOwners TicketCommentOwners
	ReportObjectTypes   ReportObjectTypes
	ReportOwners        ReportOwners
	RabbitMQ            RabbitMQConstants
}

type Context struct {
	Translator                   string
	IsLoadedValidationTranslator string
	ID                           string
	WebsocketConnection          string
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
	Email               string
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
	Room                string
	NotificationType    string
	Notification        string
	NotificationSetting string
	Panel               string
	MaintenanceRequest  string
	MaintenanceRecord   string
	Ticket              string
	Role                string
	Permission          string
	TicketComment       string
	Report              string
	ContactInformation  string
	PaymentTerm         string
	Guarantee           string
	GuaranteeViolation  string
	News                string
	Media               string
	Blog                string
	Post                string
	Like                string
	CorporationReview   string
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
	NotActive              string
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
	AlreadyBlocked         string
	AlreadyActive          string
	AlreadyResolved        string
	AlreadyArchived        string
	StatusNotChange        string
	AlreadyCanceled        string
	AlreadyRejected        string
	AlreadyAccepted        string
	AlreadyDraft           string
}

type SMSTemplates struct {
	OTP string
}

type EmailTemplates struct {
	Path            string
	PersianFileName string
	EnglishFileName string
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

type TicketOwners struct {
	User        string
	Corporation string
}

type TicketCommentOwners struct {
	User        string
	Corporation string
	Admin       string
}

type ReportObjectTypes struct {
	Maintenance string
	Panel       string
}

type ReportOwners struct {
	User string
}

type RabbitMQConstants struct {
	Exchange Exchanges
	Queue    Queues
	Headers  Headers
	Events   Events
}

type Exchanges struct {
	Notifications string
	DLX           string
	TypeTopic     string
	TypeFanout    string
}

type Queues struct {
	DLQ string
}

type Headers struct {
	RetryCount string
	LastError  string
	DeadLetter string
}

type Events struct {
	NotificationsEmail string
	NotificationsPush  string
	UserRegistered     string
	SendNotification   string
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator:                   "translator",
			IsLoadedValidationTranslator: "isLoadedValidationTranslator",
			ID:                           "ID",
			WebsocketConnection:          "wsConnection",
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
			Email:               "email",
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
			Room:                "room",
			NotificationType:    "notificationType",
			Notification:        "notification",
			Panel:               "panel",
			MaintenanceRequest:  "maintenanceRequest",
			MaintenanceRecord:   "maintenanceRecord",
			Ticket:              "ticket",
			Role:                "role",
			Permission:          "permission",
			TicketComment:       "ticketComment",
			Report:              "report",
			ContactInformation:  "contactInformation",
			NotificationSetting: "notificationSetting",
			PaymentTerm:         "paymentTerm",
			Guarantee:           "guarantee",
			GuaranteeViolation:  "guaranteeViolation",
			News:                "news",
			Media:               "media",
			Blog:                "blog",
			Post:                "post",
			Like:                "like",
			CorporationReview:   "corporationReview",
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
			NotActive:              "notActive",
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
			AlreadyBlocked:         "alreadyBlocked",
			AlreadyActive:          "alreadyActive",
			AlreadyResolved:        "alreadyResolved",
			AlreadyArchived:        "alreadyArchived",
			StatusNotChange:        "statusNotChange",
			AlreadyCanceled:        "alreadyCanceled",
			AlreadyRejected:        "alreadyRejected",
			AlreadyAccepted:        "alreadyAccepted",
			AlreadyDraft:           "alreadyDraft",
		},
		SMSTemplates: SMSTemplates{
			OTP: "sendOTPTemplate",
		},
		JWTKeysPath: JWTKeysPath{
			PublicKey:  "./internal/infrastructure/jwt/publicKey.pem",
			PrivateKey: "./internal/infrastructure/jwt/privateKey.pem",
		},
		EmailTemplates: EmailTemplates{
			Path:            "./internal/infrastructure/communication/email/templates/",
			PersianFileName: "fa.html",
			EnglishFileName: "en.html",
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
		TicketOwners: TicketOwners{
			User:        "users",
			Corporation: "corporations",
		},
		TicketCommentOwners: TicketCommentOwners{
			User:        "users",
			Corporation: "corporations",
			Admin:       "admins",
		},

		ReportObjectTypes: ReportObjectTypes{
			Maintenance: "maintenance",
			Panel:       "panel",
		},
		ReportOwners: ReportOwners{
			User: "users",
		},
		RabbitMQ: RabbitMQConstants{
			Exchange: Exchanges{
				Notifications: "notifications",
				DLX:           "dlx_notifications",
				TypeTopic:     "topic",
				TypeFanout:    "fanout",
			},
			Queue: Queues{
				DLQ: "dlq_notifications",
			},
			Headers: Headers{
				RetryCount: "x-retry-count",
				LastError:  "x-last-error",
				DeadLetter: "x-dead-letter-exchange",
			},
			Events: Events{
				NotificationsEmail: "Notifications.Email",
				NotificationsPush:  "Notifications.Push",
				UserRegistered:     "Users.Register",
				SendNotification:   "Notifications.Send",
			},
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

func (path *BucketPath) GetUserProfilePath(userID uint, pictureFileName string) string {
	return fmt.Sprintf("user/%d/profile/%s", userID, pictureFileName)
}

func (path *BucketPath) GetTicketImagePath(ticketID uint, imageFilename string) string {
	return fmt.Sprintf("tickets/%d/%s", ticketID, imageFilename)
}

func (path *BucketPath) GetCorporationLogoPath(corporationID uint, logoFileName string) string {
	return fmt.Sprintf("corporation/%d/logo/%s", corporationID, logoFileName)
}

func (path *BucketPath) GetNewsMediaPath(newsID uint, MediaFileName string) string {
	return fmt.Sprintf("news/%d/media/%s", newsID, MediaFileName)
}

func (path *BucketPath) GetNewsCoverImagePath(newsID uint, mediaFileName string) string {
	return fmt.Sprintf("news/%d/cover-image/%s", newsID, mediaFileName)
}

func (path *BucketPath) GetBlogMediaPath(corporationID uint, mediaFileName string) string {
	return fmt.Sprintf("corporation/%d/blog/media/%s", corporationID, mediaFileName)
}

func (path *BucketPath) GetBlogCoverImagePath(corporationID uint, mediaFileName string) string {
	return fmt.Sprintf("corporation/%d/blog/cover-image/%s", corporationID, mediaFileName)
}
