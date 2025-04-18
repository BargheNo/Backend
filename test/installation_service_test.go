package test

import (
	"testing"

	"github.com/BargheNo/Backend/bootstrap"
	addressdto "github.com/BargheNo/Backend/internal/application/dto/address"
	installationdto "github.com/BargheNo/Backend/internal/application/dto/installation"
	serviceimpl "github.com/BargheNo/Backend/internal/application/service"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	mocks "github.com/BargheNo/Backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	t.Run("Returns Empty List When No Requests Found", func(t *testing.T) {
		ownerID := uint(456)

		mockRequests := []*entity.InstallationRequest{}

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

		result := installationService.GetOwnerInstallationRequests(installationdto.InstallationListRequest{
			OwnerID: ownerID,
			Limit:   10,
			Offset:  0,
		})

		assert.Empty(t, result)
		assert.Len(t, result, 0)
		repo.AssertExpectations(t)
	})

	t.Run("Handles Multiple Request Statuses", func(t *testing.T) {
		ownerID := uint(101)

		mockRequests := []*entity.InstallationRequest{
			{
				Model:        database.Model{ID: 4},
				Name:         "Active Request",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      ownerID,
				BuildingType: "Residential",
			},
			{
				Model:        database.Model{ID: 5},
				Name:         "Cancelled Request",
				Status:       enum.InstallationRequestStatusCancelled,
				OwnerID:      ownerID,
				BuildingType: "Commercial",
			},
			{
				Model:        database.Model{ID: 6},
				Name:         "Expired Request",
				Status:       enum.InstallationRequestStatusExpired,
				OwnerID:      ownerID,
				BuildingType: "Industrial",
			},
		}

		mockAddress := addressdto.AddressResponse{
			Province:      "Test Province",
			City:          "Test City",
			StreetAddress: "Test Street",
			PostalCode:    "12345",
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
			ownerID,
			constants.AddressOwners.InstallationRequest,
		).Return(mockAddress).Times(3)

		result := installationService.GetOwnerInstallationRequests(installationdto.InstallationListRequest{
			OwnerID: ownerID,
			Limit:   10,
			Offset:  0,
		})

		assert.Len(t, result, 3)
		assert.Equal(t, "active", result[0].Status)
		assert.Equal(t, "cancelled", result[1].Status)
		assert.Equal(t, "expired", result[2].Status)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})
}
