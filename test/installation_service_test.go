package test

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInstallationRequestModel(t *testing.T) {
	repo := mocks.NewInstallationRepositoryMock()
	addressService := mocks.NewAddressServiceMock()
	userService := mocks.NewUserServiceMock()
	corporationService := mocks.NewCorporationServiceMock()
	chatService := mocks.NewChatServiceMock()
	db := mocks.NewDatabaseMock()
	config := bootstrap.Run()
	constants := config.Constants
	installationService := serviceimpl.NewInstallationService(
		constants,
		addressService,
		userService,
		corporationService,
		chatService,
		repo,
		db,
	)

	t.Run("Success - Get Installation Request Model", func(t *testing.T) {
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

		repo.On("FindRequestByID", db, requestID).Return(mockRequest, true).Once()

		result := installationService.GetInstallationRequestModel(requestID)

		assert.Equal(t, requestID, result.ID)
		assert.Equal(t, "Test Request", result.Name)
		assert.Equal(t, enum.InstallationRequestStatusActive, result.Status)
		assert.Equal(t, uint(1), result.Address.ID)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})

	t.Run("Error - Request Not Found", func(t *testing.T) {
		requestID := uint(456)

		repo.On("FindRequestByID", db, requestID).Return(nil, false).Once()

		assert.Panics(t, func() {
			installationService.GetInstallationRequestModel(requestID)
		})

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})
}

func TestCreateInstallationRequest(t *testing.T) {
	repo := mocks.NewInstallationRepositoryMock()
	addressService := mocks.NewAddressServiceMock()
	userService := mocks.NewUserServiceMock()
	db := mocks.NewDatabaseMock()
	config := bootstrap.Run()
	constants := config.Constants
	installationService := serviceimpl.NewInstallationService(
		constants,
		addressService,
		userService,
		nil,
		nil,
		repo,
		db,
	)

	t.Run("Success - Create Installation Request", func(t *testing.T) {
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
				OwnerType:     constants.AddressOwners.InstallationRequest,
			},
		}

		var nilRequest *entity.InstallationRequest = nil

		repo.On("FindOwnerRequestByName",
			db, requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			requestInfo.Name,
		).Return(nilRequest, false).Once()

		repo.On("FindOwnerRequests",
			db,
			requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			mock.Anything,
			mock.Anything,
		).Return([]*entity.InstallationRequest{nil}).Once()

		repo.On("CreateRequest",
			db,
			mock.MatchedBy(func(r *entity.InstallationRequest) bool {
				return r.Name == "Test Request" &&
					r.OwnerID == 123 &&
					r.BuildingType == "Residential" &&
					r.Address.ProvinceID == 1 &&
					r.Address.CityID == 2
			}),
		).Return(nil).Once()

		installationService.CreateInstallationRequest(requestInfo)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})

	t.Run("Error - Request Already Exists", func(t *testing.T) {
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
				OwnerType:     constants.AddressOwners.InstallationRequest,
			},
		}

		existingRequest := &entity.InstallationRequest{
			Name:         "Test Request",
			Status:       enum.InstallationRequestStatusActive,
			OwnerID:      123,
			BuildingType: "Residential",
		}

		repo.On("FindOwnerRequestByName",
			db, requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			requestInfo.Name,
		).Return(existingRequest, true).Once()

		assert.Panics(t, func() {
			installationService.CreateInstallationRequest(requestInfo)
		})

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})

	t.Run("Error - Too Many Active Requests", func(t *testing.T) {
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
				OwnerType:     constants.AddressOwners.InstallationRequest,
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

		repo.On("FindOwnerRequestByName",
			db, requestInfo.OwnerID,
			[]enum.InstallationRequestStatus{enum.InstallationRequestStatusActive},
			requestInfo.Name,
		).Return(nilRequest, false).Once()

		repo.On("FindOwnerRequests",
			db, requestInfo.OwnerID,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(existingRequests).Once()

		assert.Panics(t, func() {
			installationService.CreateInstallationRequest(requestInfo)
		})

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})
}

func TestGetOwnerInstallationRequests(t *testing.T) {
	repo := mocks.NewInstallationRepositoryMock()
	addressService := mocks.NewAddressServiceMock()
	userService := mocks.NewUserServiceMock()
	db := mocks.NewDatabaseMock()
	config := bootstrap.Run()
	constants := config.Constants
	installationService := serviceimpl.NewInstallationService(
		constants,
		addressService,
		userService,
		nil,
		nil,
		repo,
		db,
	)

	t.Run("Success - Get Owner Installation Requests", func(t *testing.T) {
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

		repo.On("FindOwnerRequests",
			db,
			ownerID,
			[]enum.InstallationRequestStatus{
				enum.InstallationRequestStatusActive,
				enum.InstallationRequestStatusCancelled,
				enum.InstallationRequestStatusExpired,
			},
			mock.Anything,
			mock.Anything,
		).Return(mockRequests).Once()

		addressService.On("GetAddress",
			uint(1),
			constants.AddressOwners.InstallationRequest,
		).Return(addressdto.AddressResponse{}).Once()

		addressService.On("GetAddress",
			uint(2),
			constants.AddressOwners.InstallationRequest,
		).Return(addressdto.AddressResponse{}).Once()

		result := installationService.GetOwnerInstallationRequests(installationdto.InstallationListRequest{
			OwnerID: ownerID,
			Limit:   10,
			Offset:  0,
		})

		assert.Len(t, result, 2)
		assert.Equal(t, "active", result[0].Status)
		assert.Equal(t, "cancelled", result[1].Status)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})

}

func TestGetInstallationRequest(t *testing.T) {
	repo := mocks.NewInstallationRepositoryMock()
	addressService := mocks.NewAddressServiceMock()
	userService := mocks.NewUserServiceMock()
	db := mocks.NewDatabaseMock()
	config := bootstrap.Run()
	constants := config.Constants
	installationService := serviceimpl.NewInstallationService(
		constants,
		addressService,
		userService,
		nil,
		nil,
		repo,
		db,
	)

	t.Run("Success - Get Installation Request", func(t *testing.T) {
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

		repo.On("FindRequestByID", db, requestID).Return(mockRequest, true).Once()

		addressService.On("GetAddress", requestID, constants.AddressOwners.InstallationRequest).Return(mockAddress).Once()

		userService.On("GetUserCredential", mockRequest.OwnerID).Return(userdto.CredentialResponse{
			ID: mockRequest.OwnerID,
		}).Once()

		result := installationService.GetInstallationRequest(requestID)

		assert.Equal(t, requestID, result.ID)
		assert.Equal(t, "Test Request", result.Name)
		assert.Equal(t, "active", result.Status)
		assert.Equal(t, uint(1), result.Address.ID)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})
}

func TestAddPanel(t *testing.T) {
	repo := mocks.NewInstallationRepositoryMock()
	addressService := mocks.NewAddressServiceMock()
	userService := mocks.NewUserServiceMock()
	corporationService := mocks.NewCorporationServiceMock()
	chatService := mocks.NewChatServiceMock()
	db := mocks.NewDatabaseMock()
	config := bootstrap.Run()
	constants := config.Constants
	installationService := serviceimpl.NewInstallationService(
		constants,
		addressService,
		userService,
		corporationService,
		chatService,
		repo,
		db,
	)

	t.Run("Success - Add Panel", func(t *testing.T) {
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

		corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		userService.On("FindUserByPhone",
			panelInfo.CustomerPhone,
		).Return(userdto.UserResponse{ID: customerID}).Once()

		var nilPanel *entity.Panel = nil
		repo.On("FindPanelByNameAndCustomerID",
			db,
			panelInfo.PanelName,
			customerID,
		).Return(nilPanel, false).Once()

		repo.On("CreatePanel",
			db,
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

		chatService.On("CreateOrGetRoom",
			request,
		).Return(chatdto.ChatRoomDetailsResponse{}).Once()

		installationService.AddPanel(panelInfo)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
		userService.AssertExpectations(t)
		corporationService.AssertExpectations(t)
	})

	t.Run("Error - Panel Already Exists", func(t *testing.T) {
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

		corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		userService.On("FindUserByPhone",
			panelInfo.CustomerPhone,
		).Return(userdto.UserResponse{ID: customerID}).Once()

		existingPanel := &entity.Panel{
			Name:       "Test Panel",
			CustomerID: customerID,
		}

		repo.On("FindPanelByNameAndCustomerID",
			db,
			panelInfo.PanelName,
			customerID,
		).Return(existingPanel, true).Once()

		assert.Panics(t, func() {
			installationService.AddPanel(panelInfo)
		})

		repo.AssertExpectations(t)
		userService.AssertExpectations(t)
		corporationService.AssertExpectations(t)
	})

	t.Run("Error - Invalid Corporation Access", func(t *testing.T) {
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

		corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		assert.Panics(t, func() {
			installationService.AddPanel(panelInfo)
		})

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
		userService.AssertExpectations(t)
		corporationService.AssertExpectations(t)
	})

	t.Run("Error - User Not Found", func(t *testing.T) {
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

		corporationService.On("CheckApplicantAccess",
			corporationID,
			operatorID,
		).Return(nil).Once()

		userService.On("FindUserByPhone",
			panelInfo.CustomerPhone,
		).Return(userdto.UserResponse{}).Once()

		assert.Panics(t, func() {
			installationService.AddPanel(panelInfo)
		})

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
		userService.AssertExpectations(t)
		corporationService.AssertExpectations(t)
	})
}
