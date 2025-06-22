package service

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	"github.com/BargheNo/Backend/internal/application/port"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	postgresImpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type BlogService struct {
	userService        port.UserService
	corporationService port.CorporationService
	blogRepository     postgres.BlogRepository
	constants          *bootstrap.Constants
	s3Storage          s3.S3Storage
	db                 database.Database
}

func NewBlogService(
	userService port.UserService,
	corporationService port.CorporationService,
	blogRepository postgres.BlogRepository,
	constants *bootstrap.Constants,
	s3Storage s3.S3Storage,
	db database.Database,
) *BlogService {
	return &BlogService{
		userService:        userService,
		corporationService: corporationService,
		blogRepository:     blogRepository,
		constants:          constants,
		s3Storage:          s3Storage,
		db:                 db,
	}
}

func (blogService *BlogService) mapToFilterStatuses(enumStatus uint) []enum.PostStatus {
	statuses := enum.GetAllPostStatus()
	for _, status := range statuses {
		if uint(status) == enumStatus {
			if status == enum.PostStatusAll {
				return statuses
			}
			return []enum.PostStatus{status}
		}
	}
	return statuses
}

func (blogService *BlogService) mapToOperationalStatuses(enumStatus uint) enum.PostStatus {
	allowedStatuses := []enum.PostStatus{enum.PostStatusPublished, enum.PostStatusDraft}
	for _, status := range allowedStatuses {
		if uint(status) == enumStatus {
			return status
		}
	}
	return enum.PostStatusDraft
}

func (blogService *BlogService) checkDuplicateBlog(corporationID uint, title string) error {
	post, err := blogService.blogRepository.FindCorporationPostByTitle(blogService.db, corporationID, title)
	if err != nil {
		return err
	}
	if post != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(blogService.constants.Field.Blog, blogService.constants.Tag.AlreadyExist)
		return conflictErrors
	}
	return err
}

func (blogService *BlogService) getPost(postID uint) (*entity.Post, error) {
	post, err := blogService.blogRepository.FindPostByID(blogService.db, postID)
	if err != nil {
		return nil, err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return nil, notFoundError
	}
	return post, nil
}

func (blogService *BlogService) getPostMediaByID(mediaID, postID uint) (*entity.Media, error) {
	media, err := blogService.blogRepository.FindPostMediaByID(blogService.db, mediaID, postID, "blog")
	if err != nil {
		return nil, err
	}
	if media == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Media}
		return nil, notFoundError
	}
	return media, nil
}

func (blogService *BlogService) CreatePost(request blogdto.CreatePostRequest) (uint, error) {
	if err := blogService.userService.IsUserActive(request.AuthorID); err != nil {
		return 0, err
	}

	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID); err != nil {
		return 0, err
	}

	if err := blogService.corporationService.ISCorporationApproved(request.CorporationID); err != nil {
		return 0, err
	}

	if err := blogService.checkDuplicateBlog(request.CorporationID, request.Title); err != nil {
		return 0, err
	}

	post := &entity.Post{
		Title:         request.Title,
		Content:       request.Content,
		Description:   request.Description,
		AuthorID:      request.AuthorID,
		CorporationID: request.CorporationID,
		Status:        request.Status,
	}
	err := blogService.db.WithTransaction(func(tx database.Database) error {
		if err := blogService.blogRepository.CreatePost(tx, post); err != nil {
			return err
		}

		if request.CoverImage != nil {
			post.CoverImage = blogService.constants.S3BucketPath.GetBlogCoverImagePath(request.CorporationID, request.CoverImage.Filename)
			if err := blogService.s3Storage.UploadObject(enum.BlogMedia, post.CoverImage, request.CoverImage); err != nil {
				return err
			}
		}

		if err := blogService.blogRepository.UpdatePost(tx, post); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return post.ID, nil
}

func (blogService *BlogService) GetCorporationPosts(request blogdto.GetCorporationPostsRequest) ([]blogdto.CorporationPostResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
		return nil, err
	}

	allowedStatuses := blogService.mapToFilterStatuses(request.Status)
	posts, err := blogService.blogRepository.FindCorporationPostsByStatus(blogService.db, request.CorporationID, allowedStatuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	response := make([]blogdto.CorporationPostResponse, len(posts))

	for i, post := range posts {
		coverImage := ""
		if post.CoverImage != "" {
			coverImage, err = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
			if err != nil {
				return nil, err
			}
		}

		likeCount, err := blogService.blogRepository.FindLikeCountByOwner(blogService.db, post.ID, "blog")
		if err != nil {
			return nil, err
		}

		author, err := blogService.userService.GetUserCredential(post.AuthorID)
		if err != nil {
			return nil, err
		}

		response[i] = blogdto.CorporationPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			Status:      post.Status.String(),
			Content:     post.Content,
			Author:      author,
			CoverImage:  coverImage,
			CreatedAt:   post.CreatedAt,
			LikeCount:   likeCount,
		}
	}
	return response, nil
}

func (blogService *BlogService) GetCorporationPostsForGeneral(request blogdto.GetPublicCorporationPostsRequest) ([]blogdto.GeneralPostResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	if err := blogService.corporationService.DoesCorporationExist(request.CorporationID); err != nil {
		return nil, err
	}

	allowedStatuses := []enum.PostStatus{enum.PostStatusPublished}
	posts, err := blogService.blogRepository.FindCorporationPostsByStatus(blogService.db, request.CorporationID, allowedStatuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}
	response := make([]blogdto.GeneralPostResponse, len(posts))

	for i, post := range posts {
		coverImage := ""
		if post.CoverImage != "" {
			coverImage, err = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
			if err != nil {
				return nil, err
			}
		}

		corporation, err := blogService.corporationService.GetCorporationCredentials(post.CorporationID)
		if err != nil {
			return nil, err
		}

		likeCount, err := blogService.blogRepository.FindLikeCountByOwner(blogService.db, post.ID, "blog")
		if err != nil {
			return nil, err
		}

		response[i] = blogdto.GeneralPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			Content:     post.Content,
			Corporation: corporation,
			CoverImage:  coverImage,
			CreatedAt:   post.CreatedAt,
			LikeCount:   likeCount,
		}
	}
	return response, nil
}

func (blogService *BlogService) GetGeneralPosts(request blogdto.GetPublicPostsRequest) ([]blogdto.GeneralPostResponse, error) {
	paginationModifier := postgresImpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := postgresImpl.NewSortingModifier("created_at", true)

	allowedStatuses := []enum.PostStatus{enum.PostStatusPublished}
	posts, err := blogService.blogRepository.FindPostsByStatus(blogService.db, allowedStatuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}

	response := make([]blogdto.GeneralPostResponse, len(posts))
	for i, post := range posts {
		corporation, err := blogService.corporationService.GetCorporationCredentials(post.CorporationID)
		if err != nil {
			return nil, err
		}

		coverImage := ""
		if post.CoverImage != "" {
			coverImage, err = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
			if err != nil {
				return nil, err
			}
		}

		likeCount, err := blogService.blogRepository.FindLikeCountByOwner(blogService.db, post.ID, "blog")
		if err != nil {
			return nil, err
		}

		response[i] = blogdto.GeneralPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			Content:     post.Content,
			Corporation: corporation,
			CoverImage:  coverImage,
			CreatedAt:   post.CreatedAt,
			LikeCount:   likeCount,
		}
	}
	return response, nil
}

func (blogService *BlogService) GetCorporationPost(request blogdto.GetCorporationPostRequest) (blogdto.CorporationPostResponse, error) {
	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
		return blogdto.CorporationPostResponse{}, err
	}

	post, err := blogService.blogRepository.FindCorporationPost(blogService.db, request.PostID, request.CorporationID)
	if err != nil {
		return blogdto.CorporationPostResponse{}, err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return blogdto.CorporationPostResponse{}, notFoundError
	}

	coverImage := ""
	if post.CoverImage != "" {
		coverImage, err = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
		if err != nil {
			return blogdto.CorporationPostResponse{}, err
		}
	}

	likeCount, err := blogService.blogRepository.FindLikeCountByOwner(blogService.db, post.ID, "blog")
	if err != nil {
		return blogdto.CorporationPostResponse{}, err
	}

	author, err := blogService.userService.GetUserCredential(post.AuthorID)
	if err != nil {
		return blogdto.CorporationPostResponse{}, err
	}

	return blogdto.CorporationPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		Content:     post.Content,
		Status:      post.Status.String(),
		Author:      author,
		CoverImage:  coverImage,
		CreatedAt:   post.CreatedAt,
		LikeCount:   likeCount,
	}, nil
}

func (blogService *BlogService) GetGeneralPost(postID uint) (blogdto.GeneralPostResponse, error) {
	post, err := blogService.getPost(postID)
	if err != nil {
		return blogdto.GeneralPostResponse{}, err
	}

	if post.Status == enum.PostStatusDraft {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return blogdto.GeneralPostResponse{}, notFoundError
	}

	corporation, err := blogService.corporationService.GetCorporationCredentials(post.CorporationID)
	if err != nil {
		return blogdto.GeneralPostResponse{}, err
	}

	coverImage := ""
	if post.CoverImage != "" {
		coverImage, err = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
		if err != nil {
			return blogdto.GeneralPostResponse{}, err
		}
	}

	likeCount, err := blogService.blogRepository.FindLikeCountByOwner(blogService.db, post.ID, "blog")
	if err != nil {
		return blogdto.GeneralPostResponse{}, err
	}

	return blogdto.GeneralPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		Content:     post.Content,
		Corporation: corporation,
		CoverImage:  coverImage,
		CreatedAt:   post.CreatedAt,
		LikeCount:   likeCount,
	}, nil
}

func (blogService *BlogService) EditPost(request blogdto.EditPostRequest) error {
	if err := blogService.userService.IsUserActive(request.AuthorID); err != nil {
		return err
	}

	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID); err != nil {
		return err
	}

	post, err := blogService.getPost(request.PostID)
	if err != nil {
		return err
	}

	if request.Title != nil {
		if err := blogService.checkDuplicateBlog(request.CorporationID, *request.Title); err != nil {
			return err
		}
		post.Title = *request.Title
	}

	if request.Content != nil {
		post.Content = *request.Content
	}

	if request.Description != nil {
		post.Description = *request.Description
	}

	prevCoverPath := post.CoverImage
	if request.CoverImage != nil {
		post.CoverImage = blogService.constants.S3BucketPath.GetBlogCoverImagePath(request.CorporationID, request.CoverImage.Filename)
		blogService.s3Storage.UploadObject(enum.BlogMedia, post.CoverImage, request.CoverImage)
	}

	post.Status = blogService.mapToOperationalStatuses(request.Status)

	err = blogService.db.WithTransaction(func(tx database.Database) error {
		if err = blogService.blogRepository.UpdatePost(tx, post); err != nil {
			return err
		}

		if err := blogService.s3Storage.DeleteObject(enum.BlogMedia, prevCoverPath); err != nil {
			return err
		}

		return err
	})

	return err
}

func (blogService *BlogService) DeletePost(request blogdto.DeletePostRequest) error {
	if err := blogService.userService.IsUserActive(request.AuthorID); err != nil {
		return err
	}

	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID); err != nil {
		return err
	}

	for _, postID := range request.PostIDs {
		post, err := blogService.blogRepository.FindPostByID(blogService.db, postID)
		if err != nil {
			return err
		}
		if post == nil {
			continue
		}
		blogService.blogRepository.DeletePost(blogService.db, postID)
	}
	return nil
}

func (blogService *BlogService) AddPostMedia(request blogdto.AddPostMediaRequest) (uint, error) {
	if err := blogService.userService.IsUserActive(request.AuthorID); err != nil {
		return 0, err
	}

	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID); err != nil {
		return 0, err
	}

	if _, err := blogService.getPost(request.PostID); err != nil {
		return 0, err
	}

	mediaPath := blogService.constants.S3BucketPath.GetBlogMediaPath(request.PostID, request.Media.Filename)
	if err := blogService.s3Storage.UploadObject(enum.BlogMedia, mediaPath, request.Media); err != nil {
		return 0, err
	}

	media := &entity.Media{
		Path:      mediaPath,
		OwnerID:   request.PostID,
		OwnerType: "blog",
	}
	if err := blogService.blogRepository.CreateMedia(blogService.db, media); err != nil {
		return 0, err
	}
	return media.ID, nil
}

func (blogService *BlogService) DeletePostMedia(request blogdto.AccessPostMediaRequest) error {
	if err := blogService.userService.IsUserActive(request.UserID); err != nil {
		return err
	}

	if err := blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
		return err
	}

	if _, err := blogService.getPost(request.PostID); err != nil {
		return err
	}

	media, err := blogService.getPostMediaByID(request.MediaID, request.PostID)
	if err != nil {
		return err
	}

	mediaPath := media.Path
	err = blogService.db.WithTransaction(func(tx database.Database) error {
		if err := blogService.blogRepository.DeleteMedia(tx, request.MediaID); err != nil {
			return err
		}
		if err := blogService.s3Storage.DeleteObject(enum.BlogMedia, mediaPath); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (blogService *BlogService) GetPostMedia(request blogdto.AccessPostMediaRequest) (string, error) {
	post, err := blogService.getPost(request.PostID)
	if err != nil {
		return "", err
	}

	if request.UserType == enum.UserTypeGuest && post.Status == enum.PostStatusDraft {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Media}
		return "", notFoundError
	}

	if request.UserType == enum.UserTypeCorporation {
		if err = blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID); err != nil {
			return "", err
		}
	}

	media, err := blogService.getPostMediaByID(request.MediaID, request.PostID)
	if err != nil {
		return "", err
	}

	presignedURL, err := blogService.s3Storage.GetPresignedURL(enum.BlogMedia, media.Path, 8*time.Hour)
	if err != nil {
		return "", err
	}
	return presignedURL, nil
}

func (blogService *BlogService) LikePost(request blogdto.LikePostRequest) error {
	post, err := blogService.getPost(request.PostID)
	if err != nil {
		return err
	}

	if post.Status == enum.PostStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return forbiddenError
	}

	like, err := blogService.blogRepository.FindLikeByUserAndOwner(blogService.db, request.UserID, request.PostID, "blog")
	if err != nil {
		return err
	}
	if like != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(blogService.constants.Field.Like, blogService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	like = &entity.Like{
		UserID:    request.UserID,
		OwnerID:   request.PostID,
		OwnerType: "blog",
	}
	if err := blogService.blogRepository.CreateLike(blogService.db, like); err != nil {
		return err
	}
	return nil
}

func (blogService *BlogService) UnlikePost(request blogdto.LikePostRequest) error {
	like, err := blogService.blogRepository.FindLikeByUserAndOwner(blogService.db, request.UserID, request.PostID, "blog")
	if err != nil {
		return err
	}
	if like == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Like}
		return notFoundError
	}

	if err := blogService.blogRepository.DeleteLike(blogService.db, like.ID); err != nil {
		return err
	}
	return nil
}
