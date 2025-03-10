//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/BargheNo/Backend/bootstrap"
	localizationimpl "github.com/BargheNo/Backend/internal/application/adapter/localization"
	loggerimpl "github.com/BargheNo/Backend/internal/application/adapter/logger"
	serviceimpl "github.com/BargheNo/Backend/internal/application/service"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/logger"
	"github.com/BargheNo/Backend/internal/domain/repository"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository"
	"github.com/BargheNo/Backend/internal/presentation/controller/v1/user"
	"github.com/BargheNo/Backend/internal/presentation/middleware"
	"github.com/google/wire"
)

var DatabaseProviderSet = wire.NewSet(
	repositoryimpl.NewUserRepository,
	wire.Bind(new(repository.UserRepository), new(*repositoryimpl.UserRepository)),
)

var ServiceProviderSet = wire.NewSet(
	serviceimpl.NewUserService,
	wire.Bind(new(service.UserService), new(*serviceimpl.UserService)),
)

var AdapterProviderSet = wire.NewSet(
	localizationimpl.NewTranslationService,
	loggerimpl.NewLogger,
	wire.Bind(new(logger.Logger), new(*loggerimpl.Logger)),
)

var GeneralControllerProviderSet = wire.NewSet(
	user.NewGeneralUserController,
	wire.Struct(new(GeneralControllers), "*"),
)

var ControllersProviderSet = wire.NewSet(
	wire.Struct(new(Controllers), "*"),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware.NewRecovery,
	middleware.NewLocalization,
	middleware.NewRateLimit,
	middleware.NewLoggerMiddleware,
	wire.Struct(new(Middlewares), "*"),
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

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	ServiceProviderSet,
	AdapterProviderSet,
	GeneralControllerProviderSet,
	ControllersProviderSet,
	MiddlewareProviderSet,
	ProvideConstants,
	ProvideLoggerConfig,
	ProvideRateLimitConfig,
)

type GeneralControllers struct {
	UserController *user.GeneralUserController
}

type Controllers struct {
	General *GeneralControllers
}

type Middlewares struct {
	Recovery     *middleware.RecoveryMiddleware
	Localization *middleware.LocalizationMiddleware
	RateLimit    *middleware.RateLimitMiddleware
	Logger       *middleware.LoggerMiddleware
}

type Application struct {
	Controllers *Controllers
	Middlewares *Middlewares
}

func NewApplication(
	controllers *Controllers,
	middlewares *Middlewares,
) *Application {
	return &Application{
		Controllers: controllers,
		Middlewares: middlewares,
	}
}

func InitializeApplication(container *bootstrap.Config, db database.Database) (*Application, error) {
	wire.Build(
		ProviderSet,
		NewApplication,
	)
	return &Application{}, nil
}
