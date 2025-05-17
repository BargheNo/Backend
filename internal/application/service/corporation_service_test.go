package serviceimpl

import (
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type CorporationServiceTestSuite struct {
	suite.Suite
	constants             *bootstrap.Constants
	userService           *mocks.UserServiceMock
	addressService        *mocks.AddressServiceMock
	s3Storage             *mocks.S3StorageMock
	corporationRepository *mocks.CorporationRepositoryMock
	db                    *mocks.DatabaseMock
	corporationService    *CorporationService
}

func (s *CorporationServiceTestSuite) SetupTest() {
	config := bootstrap.Run()
	s.constants = config.Constants
	s.userService = mocks.NewUserServiceMock()
	s.addressService = mocks.NewAddressServiceMock()
	s.s3Storage = mocks.NewS3StorageMock()
	s.corporationRepository = mocks.NewCorporationRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.corporationService = NewCorporationService(
		s.constants,
		s.userService,
		s.addressService,
		s.s3Storage,
		s.corporationRepository,
		s.db,
	)
}

func TestCorporationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CorporationServiceTestSuite))
}
