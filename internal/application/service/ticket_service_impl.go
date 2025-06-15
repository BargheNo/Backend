package serviceimpl

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	ticketdto "github.com/BargheNo/Backend/internal/application/dto/ticket"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
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

func (ticketService *TicketService) CreateCustomerTicket(requestInfo ticketdto.CreateTicketRequest) error {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)

	ticket := &entity.Ticket{
		Subject:     requestInfo.Subject,
		Description: requestInfo.Description,
		Status:      enum.TicketStatusNotAnswered,
		OwnerID:     requestInfo.OwnerID,
		OwnerType:   requestInfo.OwnerType,
	}
	err := ticketService.ticketRepository.CreateTicket(ticketService.db, ticket)
	if err != nil {
		return err
	}

	if requestInfo.Image != nil {
		imagePath := ticketService.constants.S3BucketPath.GetTicketImagePath(ticket.ID, requestInfo.Image.Filename)
		ticketService.s3Storage.UploadObject(enum.TicketImage, imagePath, requestInfo.Image)
		ticket.Image = imagePath
	}

	err = ticketService.ticketRepository.UpdateTicket(ticketService.db, ticket)
	if err != nil {
		return err
	}
	return nil
}

func (ticketService *TicketService) GetCustomerTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, error) {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
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
			ID: ticket.ID,
			Owner: userdto.CredentialResponse{
				FirstName: owner.FirstName,
				LastName:  owner.LastName,
				Phone:     owner.Phone,
			},
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

func (ticketService *TicketService) CreateCustomerTicketComment(requestInfo ticketdto.CreateTicketCommentRequest) error {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)
	ticket, err := ticketService.ticketRepository.GetTicketByID(ticketService.db, requestInfo.TicketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		notFoundError := exception.NotFoundError{Item: ticketService.constants.Field.Ticket}
		return notFoundError
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
	err = ticketService.ticketRepository.CreateTicketComment(ticketService.db, comment)
	if err != nil {
		return err
	}
	return nil
}

func (ticketService *TicketService) GetCustomerTicketComments(requestInfo ticketdto.TicketCommentListRequest) ([]ticketdto.TicketCommentResponse, error) {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)
	ticket, err := ticketService.ticketRepository.GetTicketByID(ticketService.db, requestInfo.TicketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		notFoundError := exception.NotFoundError{Item: ticketService.constants.Field.Ticket}
		return nil, notFoundError
	}

	if ticket.OwnerID != requestInfo.OwnerID {
		forbiddenError := exception.ForbiddenError{
			Resource: ticketService.constants.Field.Ticket,
			Message:  "",
		}
		return nil, forbiddenError
	}

	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	comments, err := ticketService.ticketRepository.GetTicketComments(ticketService.db, requestInfo.TicketID, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	responses := make([]ticketdto.TicketCommentResponse, len(comments))
	for i, comment := range comments {
		owner, err := ticketService.userService.GetUserCredential(comment.OwnerID)
		if err != nil {
			return nil, err
		}
		responses[i] = ticketdto.TicketCommentResponse{
			ID: comment.ID,
			Author: ticketdto.TicketCommentAuthorResponse{
				FirstName: owner.FirstName,
				LastName:  owner.LastName,
				OwnerType: comment.OwnerType,
			},
			Body: comment.Body,
		}
	}
	return responses, nil
}

func (ticketService *TicketService) CreateAdminTicketComment(requestInfo ticketdto.CreateTicketCommentRequest) error {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)
	ticket, err := ticketService.ticketRepository.GetTicketByID(ticketService.db, requestInfo.TicketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		notFoundError := exception.NotFoundError{Item: ticketService.constants.Field.Ticket}
		return notFoundError
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
	err = ticketService.ticketRepository.CreateTicketComment(ticketService.db, comment)
	if err != nil {
		return err
	}
	return nil
}

func (ticketService *TicketService) GetAdminTickets(requestInfo ticketdto.TicketListRequest) ([]ticketdto.TicketResponse, error) {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)
	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)
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
			ID: ticket.ID,
			Owner: userdto.CredentialResponse{
				FirstName: owner.FirstName,
				LastName:  owner.LastName,
				Phone:     owner.Phone,
			},
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
	ticket, err := ticketService.ticketRepository.GetTicketByID(ticketService.db, requestInfo.TicketID)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		notFoundError := exception.NotFoundError{Item: ticketService.constants.Field.Ticket}
		return nil, notFoundError
	}

	paginationModifier := repositoryimpl.NewPaginationModifier(requestInfo.Limit, requestInfo.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	comments, err := ticketService.ticketRepository.GetTicketComments(ticketService.db, requestInfo.TicketID, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	responses := make([]ticketdto.TicketCommentResponse, len(comments))
	for i, comment := range comments {
		owner, err := ticketService.userService.GetUserCredential(comment.OwnerID)
		if err != nil {
			return nil, err
		}
		responses[i] = ticketdto.TicketCommentResponse{
			ID: comment.ID,
			Author: ticketdto.TicketCommentAuthorResponse{
				FirstName: owner.FirstName,
				LastName:  owner.LastName,
				OwnerType: comment.OwnerType,
			},
			Body: comment.Body,
		}
	}
	return responses, nil
}

func (ticketService *TicketService) ResolveTicket(requestInfo ticketdto.ResolveTicketRequest) error {
	ticketService.userService.GetUserCredential(requestInfo.OwnerID)
	ticket, err := ticketService.ticketRepository.GetTicketByID(ticketService.db, requestInfo.TicketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		notFoundError := exception.NotFoundError{Item: ticketService.constants.Field.Ticket}
		return notFoundError
	}

	if ticket.Status == enum.TicketStatusResolved {
		var conflicterrors exception.ConflictErrors
		conflicterrors.Add(ticketService.constants.Field.Ticket, ticketService.constants.Tag.AlreadyResolved)
		return conflicterrors
	}

	ticket.Status = enum.TicketStatusResolved
	err = ticketService.ticketRepository.UpdateTicket(ticketService.db, ticket)
	if err != nil {
		return err
	}
	return nil
}
