package service

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"
	"github.com/BargheNo/Backend/internal/application/usecase"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	postgresImpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type TicketService struct {
	constants        *bootstrap.Constants
	userService      usecase.UserService
	ticketRepository postgres.TicketRepository
	s3Storage        s3.S3Storage
	db               database.Database
}

func NewTicketService(
	constants *bootstrap.Constants,
	ticketRepository postgres.TicketRepository,
	userService usecase.UserService,
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

func (ticketService *TicketService) getTicket(ticketID uint) (*entity.Ticket, error) {
	ticket, err := ticketService.ticketRepository.FindTicketByID(ticketService.db, ticketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		notFoundError := exception.NotFoundError{Item: ticketService.constants.Field.Ticket}
		return nil, notFoundError
	}
	return ticket, nil
}

func (ticketService *TicketService) CreateCustomerTicket(requestInfo ticketdto.CreateTicketRequest) error {
	ticket := &entity.Ticket{
		Subject:     requestInfo.Subject,
		Description: requestInfo.Description,
		Status:      enum.TicketStatusNotAnswered,
		OwnerID:     requestInfo.OwnerID,
		OwnerType:   requestInfo.OwnerType,
	}

	err := ticketService.db.WithTransaction(func(tx database.Database) error {
		if err := ticketService.ticketRepository.CreateTicket(tx, ticket); err != nil {
			return err
		}

		if requestInfo.Image != nil {
			ticket.Image = ticketService.constants.S3BucketPath.GetTicketImagePath(ticket.ID, requestInfo.Image.Filename)
			if err := ticketService.s3Storage.UploadObject(enum.TicketImage, ticket.Image, requestInfo.Image); err != nil {
				return err
			}
		}

		if err := ticketService.ticketRepository.UpdateTicket(tx, ticket); err != nil {
			return err
		}

		return nil
	})
	return err
}

func (ticketService *TicketService) GetCustomerTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	tickets, err := ticketService.ticketRepository.GetCustomerTickets(ticketService.db, requestInfo.OwnerID, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	responses := make([]ticketdto.TicketResponse, len(tickets))

	for i, ticket := range tickets {
		owner, err := ticketService.userService.GetUserCredential(ticket.OwnerID)
		if err != nil {
			return nil, err
		}

		responses[i] = ticketdto.TicketResponse{
			ID:          ticket.ID,
			Owner:       owner,
			Subject:     ticket.Subject.String(),
			Description: ticket.Description,
			Status:      ticket.Status.String(),
			CreatedAt:   ticket.CreatedAt,
		}

		if ticket.Image != "" {
			responses[i].Image, err = ticketService.s3Storage.GetPresignedURL(enum.TicketImage, ticket.Image, 24*time.Hour)
			if err != nil {
				return nil, err
			}
		}
	}
	return responses, nil
}

func (ticketService *TicketService) CreateCustomerTicketComment(requestInfo ticketdto.CreateTicketCommentRequest) error {
	ticket, err := ticketService.getTicket(requestInfo.TicketID)
	if err != nil {
		return err
	}

	if ticket.OwnerID != requestInfo.OwnerID {
		forbiddenError := exception.ForbiddenError{
			Resource: ticketService.constants.Field.Ticket,
			Message:  "",
		}
		return forbiddenError
	}

	if ticket.Status == enum.TicketStatusResolved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(ticketService.constants.Field.Ticket, ticketService.constants.Tag.AlreadyResolved)
		return conflictErrors
	}

	comment := &entity.TicketComment{
		TicketID:  requestInfo.TicketID,
		OwnerID:   requestInfo.OwnerID,
		OwnerType: requestInfo.OwnerType,
		Body:      requestInfo.Body,
	}
	if err = ticketService.ticketRepository.CreateTicketComment(ticketService.db, comment); err != nil {
		return err
	}
	return nil
}

func (ticketService *TicketService) GetCustomerTicketComments(requestInfo ticketdto.TicketCommentListRequest) ([]ticketdto.TicketCommentResponse, error) {
	ticket, err := ticketService.getTicket(requestInfo.TicketID)
	if err != nil {
		return nil, err
	}

	if ticket.OwnerID != requestInfo.OwnerID {
		forbiddenError := exception.ForbiddenError{
			Resource: ticketService.constants.Field.Ticket,
			Message:  "",
		}
		return nil, forbiddenError
	}

	paginationModifier := postgresImpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	comments, err := ticketService.ticketRepository.GetTicketComments(ticketService.db, requestInfo.TicketID, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	responses := make([]ticketdto.TicketCommentResponse, len(comments))

	for i, comment := range comments {
		author, err := ticketService.userService.GetUserCredential(comment.OwnerID)
		if err != nil {
			return nil, err
		}

		responses[i] = ticketdto.TicketCommentResponse{
			ID:         comment.ID,
			AuthorType: comment.OwnerType,
			Author:     author,
			Body:       comment.Body,
		}
	}
	return responses, nil
}

func (ticketService *TicketService) CreateAdminTicketComment(requestInfo ticketdto.CreateTicketCommentRequest) error {
	ticket, err := ticketService.getTicket(requestInfo.TicketID)
	if err != nil {
		return err
	}

	if ticket.Status == enum.TicketStatusResolved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(ticketService.constants.Field.Ticket, ticketService.constants.Tag.AlreadyResolved)
		return conflictErrors
	}

	comment := &entity.TicketComment{
		TicketID:  requestInfo.TicketID,
		OwnerID:   requestInfo.OwnerID,
		OwnerType: requestInfo.OwnerType,
		Body:      requestInfo.Body,
	}
	if err = ticketService.ticketRepository.CreateTicketComment(ticketService.db, comment); err != nil {
		return err
	}
	return nil
}

func (ticketService *TicketService) GetAdminTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	tickets, err := ticketService.ticketRepository.GetTickets(ticketService.db, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	responses := make([]ticketdto.TicketResponse, len(tickets))

	for i, ticket := range tickets {
		owner, err := ticketService.userService.GetUserCredential(ticket.OwnerID)
		if err != nil {
			return nil, err
		}

		responses[i] = ticketdto.TicketResponse{
			ID:          ticket.ID,
			Owner:       owner,
			Subject:     ticket.Subject.String(),
			Description: ticket.Description,
			Status:      ticket.Status.String(),
			CreatedAt:   ticket.CreatedAt,
		}

		if ticket.Image != "" {
			image, err := ticketService.s3Storage.GetPresignedURL(enum.TicketImage, ticket.Image, 24*time.Hour)
			if err != nil {
				return nil, err
			}
			responses[i].Image = image
		}
	}
	return responses, nil
}

func (ticketService *TicketService) GetAdminTicketComments(requestInfo ticketdto.TicketCommentListRequest) ([]ticketdto.TicketCommentResponse, error) {
	if _, err := ticketService.getTicket(requestInfo.TicketID); err != nil {
		return nil, err
	}

	paginationModifier := postgresImpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	comments, err := ticketService.ticketRepository.GetTicketComments(ticketService.db, requestInfo.TicketID, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	responses := make([]ticketdto.TicketCommentResponse, len(comments))

	for i, comment := range comments {
		author, err := ticketService.userService.GetUserCredential(comment.OwnerID)
		if err != nil {
			return nil, err
		}

		responses[i] = ticketdto.TicketCommentResponse{
			ID:         comment.ID,
			AuthorType: comment.OwnerType,
			Author:     author,
			Body:       comment.Body,
		}
	}
	return responses, nil
}

func (ticketService *TicketService) ResolveTicket(requestInfo ticketdto.ResolveTicketRequest) error {
	ticket, err := ticketService.getTicket(requestInfo.TicketID)
	if err != nil {
		return err
	}

	if ticket.Status == enum.TicketStatusResolved {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(ticketService.constants.Field.Ticket, ticketService.constants.Tag.AlreadyResolved)
		return conflictErrors
	}

	ticket.Status = enum.TicketStatusResolved
	if err = ticketService.ticketRepository.UpdateTicket(ticketService.db, ticket); err != nil {
		return err
	}
	return nil
}
