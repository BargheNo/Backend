package serviceimpl

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
)

type TicketService struct {
	constants        *bootstrap.Constants
	ticketRepository repository.TicketRepository
	userService      service.UserService
	s3Storage        s3.S3Storage
	db               database.Database
}

func NewTicketService(
	constants *bootstrap.Constants,
	ticketRepository repository.TicketRepository,
	userService service.UserService,
	s3Storage s3.S3Storage,
	db database.Database,
) *TicketService {
	return &TicketService{
		constants:        constants,
		ticketRepository: ticketRepository,
		userService:      userService,
		s3Storage:        s3Storage,
		db:               db,
	}
}

func (ticketService *TicketService) CreateTicket(requestInfo ticketdto.CreateTicketRequest) {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)

	ticket := &entity.Ticket{
		OwnerID:     requestInfo.OwnerID,
		Subject:     requestInfo.Subject,
		Description: requestInfo.Description,
	}
	err := ticketService.ticketRepository.CreateTicket(ticketService.db, ticket)
	if err != nil {
		panic(err)
	}

	if requestInfo.Image != nil {
		imagePath := ticketService.constants.S3BucketPath.GetTicketImagePath(ticket.ID, requestInfo.Image.Filename)
		ticketService.s3Storage.UploadObject(enum.TicketImage, imagePath, requestInfo.Image)
		ticket.Image = imagePath
	}

	err = ticketService.ticketRepository.UpdateTicket(ticketService.db, ticket)
	if err != nil {
		panic(err)
	}
}

func (ticketService *TicketService) GetCustomerTickets(requestInfo ticketdto.TicketListRequest) []ticketdto.TicketResponse {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)

	tickets := ticketService.ticketRepository.GetCustomerTickets(ticketService.db, requestInfo.OwnerID)
	responses := make([]ticketdto.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		responses[i] = ticketdto.TicketResponse{
			ID:          ticket.ID,
			OwnerID:     ticket.OwnerID,
			Subject:     ticket.Subject.String(),
			Description: ticket.Description,
			CreatedAt:   ticket.CreatedAt,
		}
		if ticket.Image != "" {
			image := ticketService.s3Storage.GetPresignedURL(enum.TicketImage, ticket.Image, 24*time.Hour)
			responses[i].Image = image
		}
	}

	return responses
}
