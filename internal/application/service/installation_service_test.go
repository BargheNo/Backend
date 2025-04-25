package serviceimpl_test

import (
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	chatdto "github.com/BargheNo/Backend/internal/application/dto/chat"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	serviceimpl "github.com/BargheNo/Backend/internal/application/service"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	mocks "github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type InstallationServiceTestSuite struct {
	suite.Suite
	repo                *mocks.InstallationRepositoryMock
	addressService      *mocks.AddressServiceMock
	userService         *mocks.UserServiceMock
	corporationService  *mocks.CorporationServiceMock
	chatService         *mocks.ChatServiceMock
	db                  *mocks.DatabaseMock
	constants           *bootstrap.Constants
	installationService *serviceimpl.InstallationService
}

func (s *InstallationServiceTestSuite) SetupTest() {
	s.repo = mocks.NewInstallationRepositoryMock()
	s.addressService = mocks.NewAddressServiceMock()
	s.userService = mocks.NewUserServiceMock()
	s.corporationService = mocks.NewCorporationServiceMock()
	s.chatService = mocks.NewChatServiceMock()
	s.db = mocks.NewDatabaseMock()
	config := bootstrap.Run()
	s.constants = config.Constants

	s.installationService = serviceimpl.NewInstallationService(
		s.constants,
		s.addressService,
		s.userService,
		s.corporationService,
		s.chatService,
		s.repo,
		s.db,
	)
}

func (s *InstallationServiceTestSuite) TestGetInstallationRequestModel() {
	s.Run("Success - Get Installation Request Model", func() {
		requestID := uint(123)
		mockRequest := &entity.InstallationRequest{
			Model:        database.Model{ID: requestID},
			Name:         "Test Request",
			Status:       enum.InstallationRequestStatusActive,
			OwnerID:      456,
			BuildingType: "Residential",
			Address: entity.Address{
				Model:         database.Model{ID: 1},
				ProvinceID:    1,
				CityID:        2,
				StreetAddress: "123 Test St",
				PostalCode:    "12345",
				HouseNumber:   "10A",
				Unit:          5,
			},
		}

		s.repo.On("FindRequestByID", s.db, requestID).Return(mockRequest, true).Once()

		result := s.installationService.GetInstallationRequestModel(requestID)

		s.Equal(requestID, result.ID)
		s.Equal("Test Request", result.Name)
		s.Equal(enum.InstallationRequestStatusActive, result.Status)
		s.Equal(uint(1), result.Address.ID)

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("Error - Request Not Found", func() {
		requestID := uint(456)

		s.repo.On("FindRequestByID", s.db, requestID).Return(nil, false).Once()

		s.Panics(func() {
			s.installationService.GetInstallationRequestModel(requestID)
		})

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *InstallationServiceTestSuite) TestCreateInstallationRequest() {
	s.Run("Success - Create Installation Request", func() {
		requestInfo := installationdto.NewInstallationRequest{
			Name:         "Test Request",
			OwnerID:      123,
			BuildingType: "Residential",
			Address: addressdto.CreateAddressRequest{
				ProvinceID:    1,
				CityID:        2,
				StreetAddress: "123 Test St",
				PostalCode:    "12345",
				HouseNumber:   "10A",
				Unit:          5,
				OwnerID:       123,
				OwnerType:     s.constants.AddressOwners.InstallationRequest,
			},
		}

		var nilRequest *entity.InstallationRequest = nil

		s.repo.On("FindOwnerRequestByName",
			s.db, requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			requestInfo.Name,
		).Return(nilRequest, false).Once()

		s.repo.On("FindOwnerRequests",
			s.db,
			requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			mock.Anything,
			mock.Anything,
		).Return([]*entity.InstallationRequest{nil}).Once()

		s.repo.On("CreateRequest",
			s.db,
			mock.MatchedBy(func(r *entity.InstallationRequest) bool {
				return r.Name == "Test Request" &&
					r.OwnerID == 123 &&
					r.BuildingType == "Residential" &&
					r.Address.ProvinceID == 1 &&
					r.Address.CityID == 2
			}),
		).Return(nil).Once()

		s.installationService.CreateInstallationRequest(requestInfo)

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("Error - Request Already Exists", func() {
		requestInfo := installationdto.NewInstallationRequest{
			Name:         "Test Request",
			OwnerID:      123,
			BuildingType: "Residential",
			Address: addressdto.CreateAddressRequest{
				ProvinceID:    1,
				CityID:        2,
				StreetAddress: "123 Test St",
				PostalCode:    "12345",
				HouseNumber:   "10A",
				Unit:          5,
				OwnerID:       123,
				OwnerType:     s.constants.AddressOwners.InstallationRequest,
			},
		}

		existingRequest := &entity.InstallationRequest{
			Name:         "Test Request",
			Status:       enum.InstallationRequestStatusActive,
			OwnerID:      123,
			BuildingType: "Residential",
		}

		s.repo.On("FindOwnerRequestByName",
			s.db, requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			requestInfo.Name,
		).Return(existingRequest, true).Once()

		s.Panics(func() {
			s.installationService.CreateInstallationRequest(requestInfo)
		})

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})

	s.Run("Error - Too Many Active Requests", func() {
		requestInfo := installationdto.NewInstallationRequest{
			Name:         "Test Request",
			OwnerID:      123,
			BuildingType: "Residential",
			Address: addressdto.CreateAddressRequest{
				ProvinceID:    1,
				CityID:        2,
				StreetAddress: "123 Test St",
				PostalCode:    "12345",
				HouseNumber:   "10A",
				Unit:          5,
				OwnerID:       123,
				OwnerType:     s.constants.AddressOwners.InstallationRequest,
			},
		}

		existingRequests := []*entity.InstallationRequest{
			{
				Name:         "Existing Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      123,
				BuildingType: "Residential",
			},
			{
				Name:         "Another Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      123,
				BuildingType: "Commercial",
			},
			{
				Name:         "Yet Another Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      123,
				BuildingType: "Industrial",
			},
			{
				Name:         "Final Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      123,
				BuildingType: "Agricultural",
			},
			{
				Name:         "Last Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      123,
				BuildingType: "Residential",
			},
		}

		var nilRequest *entity.InstallationRequest = nil

		s.repo.On("FindOwnerRequestByName",
			s.db, requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			requestInfo.Name,
		).Return(nilRequest, false).Once()

		s.repo.On("FindOwnerRequests",
			s.db, requestInfo.OwnerID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(existingRequests).Once()

		s.Panics(func() {
			s.installationService.CreateInstallationRequest(requestInfo)
		})

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *InstallationServiceTestSuite) TestGetOwnerInstallationRequests() {
	s.Run("Success - Get Owner Installation Requests", func() {
		ownerID := uint(456)

		mockRequests := []*entity.InstallationRequest{
			{
				Model:        database.Model{ID: 1},
				Name:         "Request 1",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      ownerID,
				BuildingType: "Residential",
			},
			{
				Model:        database.Model{ID: 2},
				Name:         "Request 2",
				Status:       enum.InstallationRequestStatusCancelled,
				OwnerID:      ownerID,
				BuildingType: "Commercial",
			},
		}

		s.repo.On("FindOwnerRequests",
			s.db,
			ownerID,
			[]enum.InstallationRequestStatus{
				enum.InstallationRequestStatusActive,
				enum.InstallationRequestStatusCancelled,
				enum.InstallationRequestStatusExpired,
			},
			mock.Anything,
			mock.Anything,
		).Return(mockRequests).Once()

		s.addressService.On("GetAddress",
			uint(1),
			s.constants.AddressOwners.InstallationRequest,
		).Return(addressdto.AddressResponse{}).Once()

		s.addressService.On("GetAddress",
			uint(2),
			s.constants.AddressOwners.InstallationRequest,
		).Return(addressdto.AddressResponse{}).Once()

		result := s.installationService.GetOwnerInstallationRequests(installationdto.InstallationListRequest{
			OwnerID: ownerID,
			Limit:   10,
			Offset:  0,
		})

		s.Len(result, 2)
		s.Equal("active", result[0].Status)
		s.Equal("cancelled", result[1].Status)

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *InstallationServiceTestSuite) TestGetInstallationRequest() {
	s.Run("Success - Get Installation Request", func() {
		requestID := uint(789)

		mockRequest := &entity.InstallationRequest{
			Model:        database.Model{ID: requestID},
			Name:         "Test Request",
			Status:       enum.InstallationRequestStatusActive,
			OwnerID:      123,
			BuildingType: "Residential",
			Address: entity.Address{
				Model:         database.Model{ID: 1},
				ProvinceID:    1,
				CityID:        2,
				StreetAddress: "123 Test St",
				PostalCode:    "12345",
				HouseNumber:   "10A",
				Unit:          5,
			},
		}

		mockAddress := addressdto.AddressResponse{
			ID:            1,
			Province:      "Test Province",
			City:          "Test City",
			StreetAddress: "Test Street",
			PostalCode:    "12345",
			HouseNumber:   "10A",
			Unit:          5,
		}

		s.repo.On("FindRequestByID", s.db, requestID).Return(mockRequest, true).Once()

		s.addressService.On("GetAddress", requestID, s.constants.AddressOwners.InstallationRequest).Return(mockAddress).Once()

		s.userService.On("GetUserCredential", mockRequest.OwnerID).Return(userdto.CredentialResponse{
			ID: mockRequest.OwnerID,
		}).Once()

		result := s.installationService.GetInstallationRequest(requestID)

		s.Equal(requestID, result.ID)
		s.Equal("Test Request", result.Name)
		s.Equal("active", result.Status)
		s.Equal(uint(1), result.Address.ID)

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

// func (s *InstallationServiceTestSuite) TestGetOwnerInstallationRequest() {
// 	s.Run("Success - Get Owner Installation Request", func() {
// 		requestID := uint(789)
// 		ownerID := uint(123)

// 		mockRequest := &entity.InstallationRequest{
// 			Model:        database.Model{ID: requestID},
// 			Name:         "Test Request",
// 			Status:       enum.InstallationRequestStatusActive,
// 			OwnerID:      ownerID,
// 			BuildingType: "Residential",
// 			Address: entity.Address{
// 				Model:         database.Model{ID: 1},
// 				ProvinceID:    1,
// 				CityID:        2,
// 				StreetAddress: "123 Test St",
// 				PostalCode:    "12345",
// 				HouseNumber:   "10A",
// 				Unit:          5,
// 			},
// 		}

// 		mockAddress := addressdto.AddressResponse{
// 			ID:            1,
// 			Province:      "Test Province",
// 			City:          "Test City",
// 			StreetAddress: "Test Street",
// 			PostalCode:    "12345",
// 			HouseNumber:   "10A",
// 			Unit:          5,
// 		}

// 		s.addressService.On("GetAddress", requestID, s.constants.AddressOwners.InstallationRequest).Return(mockAddress).Once()

// 		requestInfo := installationdto.GetOwnerRequest{
// 			OwnerID:   ownerID,
// 			RequestID: requestID,
// 		}

// 		result := s.installationService.GetOwnerInstallationRequest(requestInfo)

// 		s.Equal(requestID, result.ID)
// 		s.Equal("Test Request", result.Name)
// 		s.Equal("active", result.Status)
// 		s.Equal(uint(1), result.Address.ID)

// 		s.repo.AssertExpectations(s.T())
// 		s.addressService.AssertExpectations(s.T())
// 	})
// }

func (s *InstallationServiceTestSuite) TestGetInstallationRequests() {
	s.Run("Success - Get Installation Requests", func() {

		mockRequests := []*entity.InstallationRequest{
			{
				Model:        database.Model{ID: 123},
				Name:         "Test Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      123,
				BuildingType: "Residential",
				Address: entity.Address{
					Model:         database.Model{ID: 123},
					ProvinceID:    1,
					CityID:        2,
					StreetAddress: "123 Test St",
					PostalCode:    "12345",
					HouseNumber:   "10A",
					Unit:          5,
				},
			},
			{
				Model:        database.Model{ID: 456},
				Name:         "Another Request",
				Status:       enum.InstallationRequestStatusCancelled,
				OwnerID:      456,
				BuildingType: "Commercial",
				Address: entity.Address{
					Model:         database.Model{ID: 456},
					ProvinceID:    1,
					CityID:        2,
					StreetAddress: "456 Test St",
					PostalCode:    "67890",
					HouseNumber:   "20B",
					Unit:          10,
				},
			},
		}

		s.repo.On("FindRequestByStatus",
			s.db,
			[]enum.InstallationRequestStatus{
				enum.InstallationRequestStatusActive,
			},
		).Return(mockRequests).Once()

		s.addressService.On("GetAddress",
			uint(123),
			s.constants.AddressOwners.InstallationRequest,
		).Return(addressdto.AddressResponse{ID: 123}).Once()

		s.addressService.On("GetAddress",
			uint(456),
			s.constants.AddressOwners.InstallationRequest,
		).Return(addressdto.AddressResponse{ID: 456}).Once()

		requestInfo := installationdto.InstallationListRequest{
			OwnerID: 123,
			Offset:  0,
			Limit:   10,
		}

		result := s.installationService.GetInstallationRequests(requestInfo)

		s.Len(result, 2)
		s.Equal(123, result[0].ID)
		s.Equal(456, result[1].ID)
		s.Equal("Test Request", result[0].Name)
		s.Equal("Another Request", result[1].Name)
		s.Equal("active", result[0].Status)
		s.Equal("cancelled", result[1].Status)
		s.Equal(uint(123), result[0].Address.ID)
		s.Equal(uint(456), result[1].Address.ID)

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
	})
}

func (s *InstallationServiceTestSuite) TestAddPanel() {
	s.Run("Success - Add Panel", func() {
		operatorID := uint(456)
		corporationID := uint(123)
		customerID := uint(789)

		panelInfo := installationdto.AddPanelRequest{
			CorporationID:        corporationID,
			OperatorID:           operatorID,
			PanelName:            "Test Panel",
			CustomerPhone:        "1234567890",
			Power:                1000,
			Area:                 50,
			BuildingType:         "Residential",
			Tilt:                 30,
			Azimuth:              45,
			TotalNumberOfModules: 10,
		}

		s.corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		s.userService.On("FindUserByPhone",
			panelInfo.CustomerPhone,
		).Return(userdto.UserResponse{ID: customerID}).Once()

		var nilPanel *entity.Panel = nil
		s.repo.On("FindPanelByNameAndCustomerID",
			s.db,
			panelInfo.PanelName,
			customerID,
		).Return(nilPanel, false).Once()

		s.repo.On("CreatePanel",
			s.db,
			mock.MatchedBy(func(panel *entity.Panel) bool {
				return panel.Name == "Test Panel" &&
					panel.CustomerID == customerID &&
					panel.Power == 1000 &&
					panel.Area == 50 &&
					panel.BuildingType == "Residential" &&
					panel.Tilt == 30 &&
					panel.Azimuth == 45 &&
					panel.TotalNumberOfModules == 10
			},
			),
		).Return(nil).Once()

		request := chatdto.CreateOrGetUserRoomRequest{
			CorporationID: corporationID,
			UserID:        customerID,
		}

		s.chatService.On("CreateOrGetRoom",
			request,
		).Return(chatdto.ChatRoomDetailsResponse{}).Once()

		s.installationService.AddPanel(panelInfo)

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.corporationService.AssertExpectations(s.T())
	})

	s.Run("Error - Panel Already Exists", func() {
		operatorID := uint(456)
		corporationID := uint(123)
		customerID := uint(789)

		panelInfo := installationdto.AddPanelRequest{
			CorporationID:        corporationID,
			OperatorID:           operatorID,
			PanelName:            "Test Panel",
			CustomerPhone:        "1234567890",
			Power:                1000,
			Area:                 50,
			BuildingType:         "Residential",
			Tilt:                 30,
			Azimuth:              45,
			TotalNumberOfModules: 10,
		}

		s.corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		s.userService.On("FindUserByPhone",
			panelInfo.CustomerPhone,
		).Return(userdto.UserResponse{ID: customerID}).Once()

		existingPanel := &entity.Panel{
			Name:       "Test Panel",
			CustomerID: customerID,
		}

		s.repo.On("FindPanelByNameAndCustomerID",
			s.db,
			panelInfo.PanelName,
			customerID,
		).Return(existingPanel, true).Once()

		s.Panics(func() {
			s.installationService.AddPanel(panelInfo)
		})

		s.repo.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.corporationService.AssertExpectations(s.T())
	})

	s.Run("Error - Invalid Corporation Access", func() {
		operatorID := uint(456)
		corporationID := uint(123)

		panelInfo := installationdto.AddPanelRequest{
			CorporationID:        corporationID,
			OperatorID:           operatorID,
			PanelName:            "Test Panel",
			CustomerPhone:        "1234567890",
			Power:                1000,
			Area:                 50,
			BuildingType:         "Residential",
			Tilt:                 30,
			Azimuth:              45,
			TotalNumberOfModules: 10,
		}

		s.corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		s.Panics(func() {
			s.installationService.AddPanel(panelInfo)
		})

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.corporationService.AssertExpectations(s.T())
	})

	s.Run("Error - User Not Found", func() {
		operatorID := uint(456)
		corporationID := uint(123)

		panelInfo := installationdto.AddPanelRequest{
			CorporationID:        corporationID,
			OperatorID:           operatorID,
			PanelName:            "Test Panel",
			CustomerPhone:        "1234567890",
			Power:                1000,
			Area:                 50,
			BuildingType:         "Residential",
			Tilt:                 30,
			Azimuth:              45,
			TotalNumberOfModules: 10,
		}

		s.corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		s.userService.On("FindUserByPhone",
			panelInfo.CustomerPhone,
		).Return(userdto.UserResponse{}).Once()

		s.Panics(func() {
			s.installationService.AddPanel(panelInfo)
		})

		s.repo.AssertExpectations(s.T())
		s.addressService.AssertExpectations(s.T())
		s.userService.AssertExpectations(s.T())
		s.corporationService.AssertExpectations(s.T())
	})
}

func TestInstallationService(t *testing.T) {
	suite.Run(t, new(InstallationServiceTestSuite))
}
