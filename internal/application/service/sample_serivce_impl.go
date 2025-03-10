package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/repository"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type SampleService struct {
	constants        *bootstrap.Constants
	sampleRepository repository.SampleRepository
	db               database.Database
}

func NewSampleService(
	constants *bootstrap.Constants,
	sampleRepository repository.SampleRepository,
	db database.Database,
) *SampleService {
	return &SampleService{
		constants:        constants,
		sampleRepository: sampleRepository,
		db:               db,
	}
}

func (sampleService *SampleService) SampleCreate() {
	user := &entity.User{
		FirstName:    "ali",
		LastName:     "ali zadeh",
		Phone:        "1111",
		Email:        "daryoosh.safe@gmail.com",
		Password:     "aaa",
		NationalCode: "123321",
	}
	if err := sampleService.sampleRepository.Create(sampleService.db, user); err != nil {
		panic(err)
	}
}

func (sampleService *SampleService) SampleDelete() {
	if err := sampleService.sampleRepository.Delete(sampleService.db, 1); err != nil {
		panic(err)
	}
}
