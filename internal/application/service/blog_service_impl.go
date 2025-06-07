package serviceimpl

import (
	"time"

	"github.com/BargheNo/Backend/bootstrap"
	loggerimpl "github.com/BargheNo/Backend/internal/application/adapter/logger"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/BargheNo/Backend/internal/domain/logger"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
	repositoryimpl "github.com/BargheNo/Backend/internal/infrastructure/repository/postgres"
)

type BlogService struct {
	userService        service.UserService
	corporationService service.CorporationService
	blogRepository     repository.BlogRepository
	constants          *bootstrap.Constants
	s3Storage          s3.S3Storage
	db                 database.Database
}

func NewBlogService(
	userService service.UserService,
	corporationService service.CorporationService,
	blogRepository repository.BlogRepository,
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

func (blogService *BlogService) CreatePost(request blogdto.CreatePostRequest) error {
	ok, err := blogService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return forbiddenError
	}
	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	post := &entity.Post{
		Title:         request.Title,
		Content:       request.Content,
		Description:   request.Description,
		AuthorID:      request.AuthorID,
		CorporationID: request.CorporationID,
		Status:        enum.PostStatus(request.Status),
	}

	if err := blogService.blogRepository.CreatePost(blogService.db, post); err != nil {
		return err
	}

	if request.CoverImage != nil {
		coverImagePath := blogService.constants.S3BucketPath.GetBlogCoverImagePath(request.CorporationID, request.CoverImage.Filename)
		blogService.s3Storage.UploadObject(enum.BlogMedia, coverImagePath, request.CoverImage)
		post.CoverImage = coverImagePath
	}

	err = blogService.blogRepository.UpdatePost(blogService.db, post)
	if err != nil {
		return err
	}
	return nil
}

func (blogService *BlogService) GetCorporationPosts(request blogdto.GetPostsRequest) ([]blogdto.CorporationPostResponse, error) {
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	ok, err := blogService.userService.IsUserActive(request.UserID)
	if err != nil {
		return nil, err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return nil, &forbiddenError
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)

	posts, err := blogService.blogRepository.GetCorporationPostsByStatus(blogService.db, request.CorporationID, request.Statuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}

	response := make([]blogdto.CorporationPostResponse, len(posts))
	for i, post := range posts {
		coverImage := ""
		if post.CoverImage != "" {
			coverImage = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
		}
		likeCount := blogService.blogRepository.GetLikeCountByOwner(blogService.db, post.ID, "blog")

		author, err := blogService.userService.GetUserCredential(post.AuthorID)
		if err != nil {
			return nil, err
		}

		response[i] = blogdto.CorporationPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			Status:      uint(post.Status),
			Content:     post.Content,
			Author:      author.FirstName + " " + author.LastName,
			CoverImage:  coverImage,
			CreatedAt:   post.CreatedAt,
			LikeCount:   likeCount,
		}
	}
	return response, nil
}

func (blogService *BlogService) GetCorporationPostsForGeneral(request blogdto.GetPostsRequest) ([]blogdto.GeneralPostResponse, error) {
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	blogService.corporationService.DoesCorporationExist(request.CorporationID)

	posts, err := blogService.blogRepository.GetCorporationPostsByStatus(blogService.db, request.CorporationID, request.Statuses, paginationModifier, sortingModifier)
	if err != nil {
		return nil, err
	}

	response := make([]blogdto.GeneralPostResponse, len(posts))
	for i, post := range posts {
		coverImage := ""
		if post.CoverImage != "" {
			coverImage = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
		}

		corporation, err := blogService.corporationService.GetCorporationCredentials(post.CorporationID)
		if err != nil {
			return nil, err
		}

		likeCount := blogService.blogRepository.GetLikeCountByOwner(blogService.db, post.ID, "blog")

		response[i] = blogdto.GeneralPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			Status:      uint(post.Status),
			Content:     post.Content,
			Corporation: corporation,
			CoverImage:  coverImage,
			CreatedAt:   post.CreatedAt,
			LikeCount:   likeCount,
		}
	}
	return response, nil
}
func (blogService *BlogService) GetGeneralPosts(request blogdto.GetPostsRequest) ([]blogdto.GeneralPostResponse, error) {
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	posts, err := blogService.blogRepository.GetPostsByStatus(blogService.db, request.Statuses, paginationModifier, sortingModifier)
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
			coverImage = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
		}
		likeCount := blogService.blogRepository.GetLikeCountByOwner(blogService.db, post.ID, "blog")

		response[i] = blogdto.GeneralPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Description: post.Description,
			Status:      uint(post.Status),
			Content:     post.Content,
			Corporation: corporation,
			CoverImage:  coverImage,
			CreatedAt:   post.CreatedAt,
			LikeCount:   likeCount,
		}
	}
	return response, nil
}

func (blogService *BlogService) GetCorporationPost(request blogdto.GetPostRequest) (blogdto.CorporationPostResponse, error) {
	ok, err := blogService.userService.IsUserActive(request.UserID)
	if err != nil {
		return blogdto.CorporationPostResponse{}, err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return blogdto.CorporationPostResponse{}, &forbiddenError
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)

	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return blogdto.CorporationPostResponse{}, err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return blogdto.CorporationPostResponse{}, &notFoundError
	}

	coverImage := ""
	if post.CoverImage != "" {
		coverImage = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
	}

	likeCount := blogService.blogRepository.GetLikeCountByOwner(blogService.db, post.ID, "blog")

	author, err := blogService.userService.GetUserCredential(post.AuthorID)
	if err != nil {
		return blogdto.CorporationPostResponse{}, err
	}

	return blogdto.CorporationPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		Content:     post.Content,
		Status:      uint(post.Status),
		Author:      author.FirstName + " " + author.LastName,
		CoverImage:  coverImage,
		CreatedAt:   post.CreatedAt,
		LikeCount:   likeCount,
	}, nil
}

func (blogService *BlogService) GetGeneralPost(request blogdto.GetPostRequest) (blogdto.GeneralPostResponse, error) {
	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return blogdto.GeneralPostResponse{}, err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return blogdto.GeneralPostResponse{}, &notFoundError
	}

	if post.Status == enum.PostStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return blogdto.GeneralPostResponse{}, &forbiddenError
	}

	corporation, err := blogService.corporationService.GetCorporationCredentials(post.CorporationID)
	if err != nil {
		return blogdto.GeneralPostResponse{}, err
	}

	coverImage := ""
	if post.CoverImage != "" {
		coverImage = blogService.s3Storage.GetPresignedURL(enum.BlogMedia, post.CoverImage, 8*time.Hour)
	}

	likeCount := blogService.blogRepository.GetLikeCountByOwner(blogService.db, post.ID, "blog")

	return blogdto.GeneralPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Description: post.Description,
		Status:      uint(post.Status),
		Content:     post.Content,
		Corporation: corporation,
		CoverImage:  coverImage,
		CreatedAt:   post.CreatedAt,
		LikeCount:   likeCount,
	}, nil
}

func (blogService *BlogService) EditPost(request blogdto.EditPostRequest) error {
	ok, err := blogService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return &forbiddenError
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return &notFoundError
	}

	if request.Title != nil {
		post.Title = *request.Title
	}
	if request.Content != nil {
		post.Content = *request.Content
	}
	if request.Description != nil {
		post.Description = *request.Description
	}
	if request.CoverImage != nil {
		coverImagePath := blogService.constants.S3BucketPath.GetBlogCoverImagePath(request.CorporationID, request.CoverImage.Filename)
		blogService.s3Storage.UploadObject(enum.BlogMedia, coverImagePath, request.CoverImage)
		err := blogService.s3Storage.DeleteObject(enum.BlogMedia, post.CoverImage)
		if err != nil {
			loggerimpl.GetLogger().Error("unable to delete object", logger.Error("error:", err))
		}
		post.CoverImage = coverImagePath
	}
	post.Status = enum.PostStatus(request.Status)

	err = blogService.blogRepository.UpdatePost(blogService.db, post)
	if err != nil {
		return err
	}
	return nil
}

func (blogService *BlogService) DeletePost(request blogdto.DeletePostRequest) error {
	ok, err := blogService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return &forbiddenError
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

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
	ok, err := blogService.userService.IsUserActive(request.AuthorID)
	if err != nil {
		return 0, err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return 0, &forbiddenError
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return 0, err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return 0, &notFoundError
	}

	mediaPath := blogService.constants.S3BucketPath.GetBlogMediaPath(request.PostID, request.Media.Filename)
	blogService.s3Storage.UploadObject(enum.BlogMedia, mediaPath, request.Media)

	media := &entity.Media{
		Path:      mediaPath,
		OwnerID:   request.PostID,
		OwnerType: "blog",
	}
	if err := blogService.blogRepository.AddMedia(blogService.db, media); err != nil {
		return 0, err
	}
	return media.ID, nil
}

func (blogService *BlogService) DeletePostMedia(request blogdto.AccessPostMediaRequest) error {
	ok, err := blogService.userService.IsUserActive(request.UserID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return &forbiddenError
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)

	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return &notFoundError
	}

	media, err := blogService.blogRepository.GetMediaByID(blogService.db, request.MediaID)
	if err != nil {
		return err
	}
	if media == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Media}
		return &notFoundError
	}

	if media.OwnerID != request.PostID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Media,
		}
		return &forbiddenError
	}

	if err := blogService.s3Storage.DeleteObject(enum.BlogMedia, media.Path); err != nil {
		return err
	}
	if err := blogService.blogRepository.DeleteMedia(blogService.db, request.MediaID); err != nil {
		return err
	}
	return nil
}

func (blogService *BlogService) GetPostMedia(request blogdto.AccessPostMediaRequest) (string, error) {
	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return "", err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return "", &notFoundError
	}

	if request.UserType == enum.UserTypeGuest && post.Status == enum.PostStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return "", &forbiddenError
	}

	if request.UserType == enum.UserTypeCorporation {
		ok, err := blogService.userService.IsUserActive(request.UserID)
		if err != nil {
			return "", err
		}
		if !ok {
			forbiddenError := exception.ForbiddenError{
				Message:  "",
				Resource: blogService.constants.Field.Post,
			}
			return "", &forbiddenError
		}
		blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)
	}

	media, err := blogService.blogRepository.GetMediaByID(blogService.db, request.MediaID)
	if err != nil {
		return "", err
	}
	if media == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Media}
		return "", &notFoundError
	}

	if media.OwnerID != request.PostID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Media,
		}
		return "", &forbiddenError
	}

	return blogService.s3Storage.GetPresignedURL(enum.BlogMedia, media.Path, 8*time.Hour), nil
}

func (blogService *BlogService) LikePost(request blogdto.LikePostRequest) error {
	ok, err := blogService.userService.IsUserActive(request.UserID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return &forbiddenError
	}

	post, err := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if err != nil {
		return err
	}
	if post == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		return &notFoundError
	}

	if post.Status == enum.PostStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return &forbiddenError
	}

	like, err := blogService.blogRepository.FindLikeByUserAndOwner(blogService.db, request.UserID, request.PostID, "blog")
	if err != nil {
		return err
	}
	if like != nil {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(blogService.constants.Field.Like, blogService.constants.Tag.AlreadyExist)
		return &conflictErrors
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
	ok, err := blogService.userService.IsUserActive(request.UserID)
	if err != nil {
		return err
	}
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		return &forbiddenError
	}

	like, err := blogService.blogRepository.FindLikeByUserAndOwner(blogService.db, request.UserID, request.PostID, "blog")
	if err != nil {
		return err
	}
	if like == nil {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Like}
		return &notFoundError
	}

	if err := blogService.blogRepository.DeleteLike(blogService.db, like.ID); err != nil {
		return err
	}
	return nil
}
