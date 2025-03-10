package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/repository"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type UserService struct {
	constants      *bootstrap.Constants
	userRepository repository.UserRepository
	db             database.Database
}

func NewUserService(
	constants *bootstrap.Constants,
	userRepository repository.UserRepository,
	db database.Database,
) *UserService {
	return &UserService{
		constants:      constants,
		userRepository: userRepository,
		db:             db,
	}
}
