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
	communicationService "github.com/BargheNo/Backend/internal/application/service/communication"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/BargheNo/Backend/internal/domain/metrics"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	cacherepository "github.com/BargheNo/Backend/internal/domain/repository/redis"
	"github.com/BargheNo/Backend/internal/domain/s3"
	cinimpl "github.com/BargheNo/Backend/internal/infrastructure/cin"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
	cacherepositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/redis"
	"github.com/BargheNo/Backend/internal/infrastructure/seed"
	"github.com/BargheNo/Backend/internal/infrastructure/storage"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/address"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/bid"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/corporation"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/installation"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/maintenance"
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
	wire.Bind(new(repository.UserRepository), new(*repositoryimpl.UserRepository)),
	wire.Bind(new(repository.InstallationRepository), new(*repositoryimpl.InstallationRepository)),
	wire.Bind(new(repository.AddressRepository), new(*repositoryimpl.AddressRepository)),
	wire.Bind(new(cacherepository.UserCacheRepository), new(*cacherepositoryimpl.UserCacheRepository)),
	wire.Bind(new(repository.CorporationRepository), new(*repositoryimpl.CorporationRepository)),
	wire.Bind(new(repository.BidRepository), new(*repositoryimpl.BidRepository)),
)

var ServiceProviderSet = wire.NewSet(
	serviceimpl.NewUserService,
	serviceimpl.NewOTPService,
	communicationService.NewSMSService,
	serviceimpl.NewJWTService,
	serviceimpl.NewInstallationService,
	serviceimpl.NewAddressService,
	serviceimpl.NewCorporationService,
	cinimpl.NewCINService,
	serviceimpl.NewBidService,
	wire.Bind(new(service.UserService), new(*serviceimpl.UserService)),
	wire.Bind(new(service.OTPService), new(*serviceimpl.OTPService)),
	wire.Bind(new(service.SMSService), new(*communicationService.SMSService)),
	wire.Bind(new(service.JWTService), new(*serviceimpl.JWTService)),
	wire.Bind(new(service.InstallationService), new(*serviceimpl.InstallationService)),
	wire.Bind(new(service.AddressService), new(*serviceimpl.AddressService)),
	wire.Bind(new(service.CorporationService), new(*serviceimpl.CorporationService)),
	wire.Bind(new(service.CINService), new(*cinimpl.CINService)),
	wire.Bind(new(service.BidService), new(*serviceimpl.BidService)),
)

var AdapterProviderSet = wire.NewSet(
	localizationimpl.NewTranslationService,
	loggerimpl.NewLogger,
	jwtimpl.NewJWTKeyManager,
	metricsimpl.NewPrometheusMetrics,
	storage.NewS3Storage,
	wire.Bind(new(logger.Logger), new(*loggerimpl.Logger)),
	wire.Bind(new(metrics.MetricsClient), new(*metricsimpl.PrometheusMetrics)),
	wire.Bind(new(s3.S3Storage), new(*storage.S3Storage)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	address.NewGeneralAddressController,
	corporation.NewGeneralCorporationController,
	wire.Struct(new(GeneralControllers), "*"),
)

var CustomerControllerProviderSet = wire.NewSet(
	user.NewCustomerUserController,
	installation.NewCustomerInstallationController,
	address.NewCustomerAddressController,
	corporation.NewCustomerCorporationController,
	bid.NewCustomerBidController,
	maintenance.NewCustomerMaintenanceController,
	wire.Struct(new(CustomerControllers), "*"),
)

var CorporationControllerProviderSet = wire.NewSet(
	corporation.NewCorporationCorporationController,
	installation.NewCorporationInstallationController,
	bid.NewCorporationBidController,
	wire.Struct(new(CorporationControllers), "*"),
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
	wire.Struct(new(Middlewares), "*"),
)

var SeederProviderSet = wire.NewSet(
	seed.NewAddressSeeder,
	wire.Struct(new(Seeds), "*"),
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

func ProvideMetrics(container *bootstrap.Config) *bootstrap.Metrics {
	return &container.Constants.Metrics
}

func ProvidePaginationConfig(container *bootstrap.Config) *bootstrap.Pagination {
	return &container.Env.Pagination
}

func ProvideStorageConfig(container *bootstrap.Config) *bootstrap.S3 {
	return &container.Env.Storage
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RepositoryProviderSet,
	ServiceProviderSet,
	AdapterProviderSet,
	GeneralControllerProviderSet,
	CustomerControllerProviderSet,
	CorporationControllerProviderSet,
	ControllersProviderSet,
	MiddlewareProviderSet,
	SeederProviderSet,
	ProvideConstants,
	ProvideLoggerConfig,
	ProvideRateLimitConfig,
	ProvideDBConfig,
	ProvideRDBConfig,
	ProvideOTPConfig,
	ProvideSMSGatewayConfig,
	ProvideSMSTemplates,
	ProvideJWTKeysPath,
	ProvideMetrics,
	ProvidePaginationConfig,
	ProvideStorageConfig,
)

type Database struct {
	DB  database.Database
	RDB database.Cache
}

type GeneralControllers struct {
	UserController        *user.GeneralUserController
	AddressController     *address.GeneralAddressController
	CorporationController *corporation.GeneralCorporationController
}

type CustomerControllers struct {
	UserController         *user.CustomerUserController
	InstallationController *installation.CustomerInstallationController
	AddressController      *address.CustomerAddressController
	CorporationController  *corporation.CustomerCorporationController
	BidController          *bid.CustomerBidController
	MaintenanceController  *maintenance.CustomerMaintenanceController
}

type CorporationControllers struct {
	CorporationController  *corporation.CorporationCorporationController
	InstallationController *installation.CorporationInstallationController
	BidController          *bid.CorporationBidController
}

type Controllers struct {
	General     *GeneralControllers
	Customer    *CustomerControllers
	Corporation *CorporationControllers
}

type Middlewares struct {
	Authentication *middleware.AuthMiddleware
	CORS           *middleware.CORSMiddleware
	Recovery       *middleware.RecoveryMiddleware
	Localization   *middleware.LocalizationMiddleware
	RateLimit      *middleware.RateLimitMiddleware
	Logger         *middleware.LoggerMiddleware
	Prometheus     *middleware.PrometheusMiddleware
}

type Seeds struct {
	AddressSeeder *seed.AddressSeeder
}

type Application struct {
	Database    *Database
	Controllers *Controllers
	Middlewares *Middlewares
	Seeds       *Seeds
}

func NewApplication(
	database *Database,
	controllers *Controllers,
	middlewares *Middlewares,
	seeds *Seeds,
) *Application {
	return &Application{
		Database:    database,
		Controllers: controllers,
		Middlewares: middlewares,
		Seeds:       seeds,
	}
}

func InitializeApplication(container *bootstrap.Config) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
