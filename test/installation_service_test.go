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
		repo,
		db,
	)

	t.Run("Success - Returns formatted installation requests", func(t *testing.T) {
		ownerID := uint(123)

		mockRequests := []*entity.InstallationRequest{
			{
				Model:        database.Model{ID: 1},
				Name:         "Request 1",
				Status:       enum.InstallationRequestStatusActive,
				OwnerID:      ownerID,
				PowerRequest: 500,
				MaxCost:      1000,
				BuildingType: "Residential",
			},
			{
				Model:        database.Model{ID: 2},
				Name:         "Request 2",
				Status:       enum.InstallationRequestStatusCancelled,
				OwnerID:      ownerID,
				PowerRequest: 750,
				MaxCost:      1500,
				BuildingType: "Commercial",
			},
		}

		mockAddress := addressdto.AddressResponse{
			Province:      "Test Province",
			City:          "Test City",
			StreetAddress: "123 Test St",
			PostalCode:    "12345",
		}

		repo.On("FindOwnerRequests",
			ownerID,
			[]enum.InstallationRequestStatus{
				enum.InstallationRequestStatusActive,
				enum.InstallationRequestStatusCancelled,
				enum.InstallationRequestStatusExpired,
			},
		).Return(mockRequests).Once()

		addressService.On("GetAddress",
			ownerID,
			constants.AddressOwners.InstallationRequest,
		).Return(mockAddress).Twice()

		result := installationService.GetOwnerInstallationRequests(installationdto.InstallationListRequest{
			OwnerID: ownerID,
			Limit:   10,
			Offset:  0,
		})

		assert.Len(t, result, 2)
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "Request 1", result[0].Name)
		assert.Equal(t, "active", result[0].Status)
		assert.Equal(t, mockAddress, result[0].Address)
		assert.Equal(t, uint(2), result[1].ID)
		assert.Equal(t, "Request 2", result[1].Name)
		assert.Equal(t, "cancelled", result[1].Status)

		repo.AssertExpectations(t)
		addressService.AssertExpectations(t)
	})

}
