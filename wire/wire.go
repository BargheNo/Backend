//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/BargheNo/Backend/bootstrap"
	jwtimpl "github.com/BargheNo/Backend/internal/application/adapter/jwt"
	localizationimpl "github.com/BargheNo/Backend/internal/application/adapter/localization"
	loggerimpl "github.com/BargheNo/Backend/internal/application/adapter/logger"
	metricsimpl "github.com/BargheNo/Backend/internal/application/adapter/metrics"
	serviceimpl "github.com/BargheNo/Backend/internal/application/service"
	"github.com/BargheNo/Backend/internal/application/service/communication/email"
	"github.com/BargheNo/Backend/internal/application/service/communication/sms"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/BargheNo/Backend/internal/domain/message"
	"github.com/BargheNo/Backend/internal/domain/metrics"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	cacherepository "github.com/BargheNo/Backend/internal/domain/repository/redis"
	"github.com/BargheNo/Backend/internal/domain/s3"
	cinimpl "github.com/BargheNo/Backend/internal/infrastructure/cin"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	"github.com/BargheNo/Backend/internal/infrastructure/rabbitmq"
	"github.com/BargheNo/Backend/internal/infrastructure/rabbitmq/consumer"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
	cacherepositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/redis"
	"github.com/BargheNo/Backend/internal/infrastructure/seed"
	"github.com/BargheNo/Backend/internal/infrastructure/storage"
	"github.com/BargheNo/Backend/internal/infrastructure/websocket"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/address"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/bid"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/blog"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/chat"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/corporation"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/installation"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/maintenance"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/news"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/notification"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/report"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/ticket"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/user"
	"github.com/BargheNo/Backend/internal/presentation/middleware"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	database.NewPostgresDatabase,
	database.NewRedisDatabase,
	wire.Bind(new(database.Database), new(*database.PostgresDatabase)),
	wire.Bind(new(database.Cache), new(*database.RedisDatabase)),
	wire.Struct(new(Database), "*"),
)

var RepositoryProviderSet = wire.NewSet(
	repositoryimpl.NewUserRepository,
	repositoryimpl.NewInstallationRepository,
	repositoryimpl.NewAddressRepository,
	cacherepositoryimpl.NewUserCacheRepository,
	repositoryimpl.NewCorporationRepository,
	repositoryimpl.NewBidRepository,
	repositoryimpl.NewChatRepository,
	repositoryimpl.NewNotificationRepository,
	repositoryimpl.NewMaintenanceRepository,
	repositoryimpl.NewTicketRepository,
	repositoryimpl.NewReportRepository,
	repositoryimpl.NewNewsRepository,
	repositoryimpl.NewBlogRepository,
	wire.Bind(new(repository.UserRepository), new(*repositoryimpl.UserRepository)),
	wire.Bind(new(repository.InstallationRepository), new(*repositoryimpl.InstallationRepository)),
	wire.Bind(new(repository.AddressRepository), new(*repositoryimpl.AddressRepository)),
	wire.Bind(new(cacherepository.UserCacheRepository), new(*cacherepositoryimpl.UserCacheRepository)),
	wire.Bind(new(repository.CorporationRepository), new(*repositoryimpl.CorporationRepository)),
	wire.Bind(new(repository.BidRepository), new(*repositoryimpl.BidRepository)),
	wire.Bind(new(repository.ChatRepository), new(*repositoryimpl.ChatRepository)),
	wire.Bind(new(repository.NotificationRepository), new(*repositoryimpl.NotificationRepository)),
	wire.Bind(new(repository.MaintenanceRepository), new(*repositoryimpl.MaintenanceRepository)),
	wire.Bind(new(repository.TicketRepository), new(*repositoryimpl.TicketRepository)),
	wire.Bind(new(repository.ReportRepository), new(*repositoryimpl.ReportRepository)),
	wire.Bind(new(repository.NewsRepository), new(*repositoryimpl.NewsRepository)),
	wire.Bind(new(repository.BlogRepository), new(*repositoryimpl.BlogRepository)),
)

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(serviceimpl.UserServiceDeps), "*"),
	wire.Struct(new(serviceimpl.NotificationServiceDeps), "*"),
	serviceimpl.NewUserService,
	serviceimpl.NewOTPService,
	sms.NewSMSService,
	email.NewEmailService,
	serviceimpl.NewJWTService,
	serviceimpl.NewInstallationService,
	serviceimpl.NewAddressService,
	serviceimpl.NewCorporationService,
	cinimpl.NewCINService,
	serviceimpl.NewBidService,
	serviceimpl.NewChatService,
	serviceimpl.NewNotificationService,
	serviceimpl.NewMaintenanceService,
	serviceimpl.NewTicketService,
	serviceimpl.NewReportService,
	serviceimpl.NewNewsService,
	serviceimpl.NewBlogService,
	wire.Bind(new(service.UserService), new(*serviceimpl.UserService)),
	wire.Bind(new(service.OTPService), new(*serviceimpl.OTPService)),
	wire.Bind(new(service.SMSService), new(*sms.SMSService)),
	wire.Bind(new(service.EmailService), new(*email.EmailService)),
	wire.Bind(new(service.JWTService), new(*serviceimpl.JWTService)),
	wire.Bind(new(service.InstallationService), new(*serviceimpl.InstallationService)),
	wire.Bind(new(service.AddressService), new(*serviceimpl.AddressService)),
	wire.Bind(new(service.CorporationService), new(*serviceimpl.CorporationService)),
	wire.Bind(new(service.CINService), new(*cinimpl.CINService)),
	wire.Bind(new(service.BidService), new(*serviceimpl.BidService)),
	wire.Bind(new(service.ChatService), new(*serviceimpl.ChatService)),
	wire.Bind(new(service.NotificationService), new(*serviceimpl.NotificationService)),
	wire.Bind(new(service.MaintenanceService), new(*serviceimpl.MaintenanceService)),
	wire.Bind(new(service.TicketService), new(*serviceimpl.TicketService)),
	wire.Bind(new(service.ReportService), new(*serviceimpl.ReportService)),
	wire.Bind(new(service.NewsService), new(*serviceimpl.NewsService)),
	wire.Bind(new(service.BlogService), new(*serviceimpl.BlogService)),
)

var AdapterProviderSet = wire.NewSet(
	localizationimpl.NewTranslationService,
	loggerimpl.NewLogger,
	jwtimpl.NewJWTKeyManager,
	metricsimpl.NewPrometheusMetrics,
	storage.NewS3Storage,
	rabbitmq.NewRabbitMQ,
	wire.Bind(new(logger.Logger), new(*loggerimpl.Logger)),
	wire.Bind(new(metrics.MetricsClient), new(*metricsimpl.PrometheusMetrics)),
	wire.Bind(new(s3.S3Storage), new(*storage.S3Storage)),
	wire.Bind(new(message.Broker), new(*rabbitmq.RabbitMQ)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	address.NewGeneralAddressController,
	corporation.NewGeneralCorporationController,
	notification.NewGeneralNotificationController,
	news.NewGeneralNewsController,
	blog.NewGeneralBlogController,
	wire.Struct(new(GeneralControllers), "*"),
)

var CustomerControllerProviderSet = wire.NewSet(
	user.NewCustomerUserController,
	installation.NewCustomerInstallationController,
	address.NewCustomerAddressController,
	corporation.NewCustomerCorporationController,
	bid.NewCustomerBidController,
	chat.NewCustomerChatController,
	notification.NewCustomerNotificationController,
	maintenance.NewCustomerMaintenanceController,
	ticket.NewCustomerTicketController,
	report.NewCustomerReportController,
	wire.Struct(new(CustomerControllers), "*"),
)

var CorporationControllerProviderSet = wire.NewSet(
	corporation.NewCorporationCorporationController,
	installation.NewCorporationInstallationController,
	chat.NewCorporationChatController,
	bid.NewCorporationBidController,
	maintenance.NewCorporationMaintenanceController,
	blog.NewCorporationBlogController,
	wire.Struct(new(CorporationControllers), "*"),
)

var AdminControllerProviderSet = wire.NewSet(
	ticket.NewAdminTicketController,
	user.NewAdminUserController,
	report.NewAdminReportController,
	news.NewAdminNewsController,
	wire.Struct(new(AdminControllers), "*"),
)

var ControllersProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewAuthMiddleware,
	middleware.NewCorsMiddleware,
	middleware.NewRecovery,
	middleware.NewLocalization,
	middleware.NewRateLimit,
	middleware.NewLoggerMiddleware,
	middleware.NewPrometheusMiddleware,
	middleware.NewWebsocketMiddleware,
	wire.Struct(new(Middlewares), "*"),
)

var SeederProviderSet = wire.NewSet(
	seed.NewAddressSeeder,
	seed.NewNotificationTypeSeeder,
	seed.NewRoleSeeder,
	seed.NewContactTypeSeeder,
	wire.Struct(new(Seeds), "*"),
)

var ConsumerProviderSet = wire.NewSet(
	consumer.NewRegisterConsumer,
	consumer.NewPushConsumer,
	consumer.NewEmailConsumer,
	consumer.NewSendNotificationConsumer,
	wire.Struct(new(Consumers), "*"),
)

func ProvideConstants(container *bootstrap.Config) *bootstrap.Constants {
	return container.Constants
}

func ProvideLoggerConfig(container *bootstrap.Config) *bootstrap.Logger {
	return &container.Env.Logger
}

func ProvideRateLimitConfig(container *bootstrap.Config) *bootstrap.RateLimit {
	return &container.Env.RateLimit
}

func ProvideDBConfig(container *bootstrap.Config) *bootstrap.Database {
	return &container.Env.PrimaryDB
}

func ProvideRDBConfig(container *bootstrap.Config) *bootstrap.Redis {
	return &container.Env.PrimaryRedis
}

func ProvideOTPConfig(container *bootstrap.Config) *bootstrap.OTP {
	return &container.Env.OTP
}

func ProvideSMSGatewayConfig(container *bootstrap.Config) *bootstrap.SMSGateway {
	return &container.Env.SMSGateway
}

func ProvideSMSTemplates(container *bootstrap.Config) *bootstrap.SMSTemplates {
	return &container.Constants.SMSTemplates
}

func ProvideJWTKeysPath(container *bootstrap.Config) *bootstrap.JWTKeysPath {
	return &container.Constants.JWTKeysPath
}

func ProvideEmailTemplates(container *bootstrap.Config) *bootstrap.EmailTemplates {
	return &container.Constants.EmailTemplates
}

func ProvideMetrics(container *bootstrap.Config) *bootstrap.Metrics {
	return &container.Constants.Metrics
}

func ProvidePaginationConfig(container *bootstrap.Config) *bootstrap.Pagination {
	return &container.Env.Pagination
}

func ProvideStorageConfig(container *bootstrap.Config) *bootstrap.S3 {
	return &container.Env.Storage
}

func ProvideWebsocketSetting(container *bootstrap.Config) *bootstrap.WebsocketSetting {
	return &container.Env.WebsocketSetting
}

func ProvideEmailSenderAccount(container *bootstrap.Config) *bootstrap.EmailAccount {
	return &container.Env.EmailSenderAccount
}

func ProvideSuperAdminCredential(container *bootstrap.Config) *bootstrap.AdminCredentials {
	return &container.Env.SuperAdmin
}

func ProvideRabbitMQConfig(container *bootstrap.Config) *bootstrap.RabbitMQ {
	return &container.Env.RabbitMQ
}

func ProvideRabbitMQConstants(container *bootstrap.Config) *bootstrap.RabbitMQConstants {
	return &container.Constants.RabbitMQ
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	AdapterProviderSet,
	GeneralControllerProviderSet,
	CustomerControllerProviderSet,
	CorporationControllerProviderSet,
	AdminControllerProviderSet,
	ControllersProviderSet,
	MiddlewareProviderSet,
	SeederProviderSet,
	ConsumerProviderSet,
	ProvideConstants,
	ProvideLoggerConfig,
	ProvideRateLimitConfig,
	ProvideDBConfig,
	ProvideRDBConfig,
	ProvideOTPConfig,
	ProvideSMSGatewayConfig,
	ProvideSMSTemplates,
	ProvideEmailTemplates,
	ProvideJWTKeysPath,
	ProvideMetrics,
	ProvidePaginationConfig,
	ProvideStorageConfig,
	ProvideWebsocketSetting,
	ProvideEmailSenderAccount,
	ProvideSuperAdminCredential,
	ProvideRabbitMQConfig,
	ProvideRabbitMQConstants,
)

type Database struct {
	DB  database.Database
	RDB database.Cache
}

type GeneralControllers struct {
	UserController         *user.GeneralUserController
	AddressController      *address.GeneralAddressController
	CorporationController  *corporation.GeneralCorporationController
	NotificationController *notification.GeneralNotificationController
	NewsController         *news.GeneralNewsController
	BlogController         *blog.GeneralBlogController
}

type CustomerControllers struct {
	UserController         *user.CustomerUserController
	InstallationController *installation.CustomerInstallationController
	AddressController      *address.CustomerAddressController
	CorporationController  *corporation.CustomerCorporationController
	BidController          *bid.CustomerBidController
	ChatController         *chat.CustomerChatController
	NotificationController *notification.CustomerNotificationController
	MaintenanceController  *maintenance.CustomerMaintenanceController
	TicketController       *ticket.CustomerTicketController
	ReportController       *report.CustomerReportController
}

type CorporationControllers struct {
	CorporationController  *corporation.CorporationCorporationController
	InstallationController *installation.CorporationInstallationController
	ChatController         *chat.CorporationChatController
	BidController          *bid.CorporationBidController
	MaintenanceController  *maintenance.CorporationMaintenanceController
	BlogController         *blog.CorporationBlogController
}

type AdminControllers struct {
	TicketController *ticket.AdminTicketController
	UserController   *user.AdminUserController
	ReportController *report.AdminReportController
	NewsController   *news.AdminNewsController
}

type Controllers struct {
	General     *GeneralControllers
	Customer    *CustomerControllers
	Corporation *CorporationControllers
	Admin       *AdminControllers
}

type Middlewares struct {
	Authentication      *middleware.AuthMiddleware
	CORS                *middleware.CORSMiddleware
	Recovery            *middleware.RecoveryMiddleware
	Localization        *middleware.LocalizationMiddleware
	RateLimit           *middleware.RateLimitMiddleware
	Logger              *middleware.LoggerMiddleware
	Prometheus          *middleware.PrometheusMiddleware
	WebsocketMiddleware *middleware.WebsocketMiddleware
}

type Seeds struct {
	AddressSeeder          *seed.AddressSeeder
	NotificationTypeSeeder *seed.NotificationTypeSeeder
	RoleSeeder             *seed.RoleSeeder
	ContactType            *seed.ContactTypeSeeder
}

type Consumers struct {
	Register     *consumer.RegisterConsumer
	Push         *consumer.PushConsumer
	Email        *consumer.EmailConsumer
	Notification *consumer.SendNotificationConsumer
}

type Application struct {
	Database    *Database
	Controllers *Controllers
	Middlewares *Middlewares
	Seeds       *Seeds
	Consumers   *Consumers
}

func NewApplication(
	database *Database,
	controllers *Controllers,
	middlewares *Middlewares,
	seeds *Seeds,
	consumers *Consumers,
) *Application {
	return &Application{
		Database:    database,
		Controllers: controllers,
		Middlewares: middlewares,
		Seeds:       seeds,
		Consumers:   consumers,
	}
}

func InitializeApplication(container *bootstrap.Config, hub *websocket.Hub) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
