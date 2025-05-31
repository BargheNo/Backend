package serviceimpl

import (
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

func (blogService *BlogService) CreatePost(request blogdto.CreatePostRequest) {
	ok := blogService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		panic(forbiddenError)
	}
	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	post := &entity.Post{
		Title:         request.Title,
		Content:       request.Content,
		AuthorID:      request.AuthorID,
		CorporationID: request.CorporationID,
		Status:        enum.PostStatus(request.Status),
	}

	if err := blogService.blogRepository.CreatePost(blogService.db, post); err != nil {
		panic(err)
	}

	if request.CoverImage != nil {
		coverImagePath := blogService.constants.S3BucketPath.GetBlogCoverImagePath(request.CorporationID, request.CoverImage.Filename)
		blogService.s3Storage.UploadObject(enum.BlogMedia, coverImagePath, request.CoverImage)
		post.CoverImage = coverImagePath
	}

	err := blogService.blogRepository.UpdatePost(blogService.db, post)
	if err != nil {
		panic(err)
	}
}

func (blogService *BlogService) GetPosts(request blogdto.GetPostsRequest) []blogdto.PostResponse {
	paginationModifier := repositoryimpl.NewPaginationModifier(request.Limit, request.Offset)
	sortingModifier := repositoryimpl.NewSortingModifier("created_at", true)

	if request.UserType == enum.UserTypeCorporation {
		ok := blogService.userService.IsUserActive(request.UserID)
		if !ok {
			forbiddenError := exception.ForbiddenError{
				Message:  "",
				Resource: blogService.constants.Field.Post,
			}
			panic(forbiddenError)
		}
		blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.UserID)
	}

	posts := blogService.blogRepository.GetCorporationPostsByStatus(blogService.db, request.CorporationID, request.Statuses, paginationModifier, sortingModifier)

	response := make([]blogdto.PostResponse, len(posts))
	for i, post := range posts {
		// corporation := blogService.corporationService.GetCorporationCredentials(post.CorporationID)
		author := blogService.userService.GetUserCredential(post.AuthorID)
		response[i] = blogdto.PostResponse{
			ID:     post.ID,
			Title:  post.Title,
			Status: uint(post.Status),
			// Corporation: corporation.Name,
			Author:     author.FirstName + " " + author.LastName,
			CoverImage: post.CoverImage,
			CreatedAt:  post.CreatedAt,
		}
	}
	return response
}

func (blogService *BlogService) GetPost(request blogdto.GetPostRequest) blogdto.PostDetailsResponse {
	post, exist := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		panic(notFoundError)
	}

	if request.UserType == enum.UserTypeCorporation {
		ok := blogService.userService.IsUserActive(request.UserID)
		if !ok {
			forbiddenError := exception.ForbiddenError{
				Message:  "",
				Resource: blogService.constants.Field.Post,
			}
			panic(forbiddenError)
		}
	}

	if request.UserType == enum.UserTypeGuest && post.Status == enum.PostStatusDraft {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		panic(forbiddenError)
	}
	author := blogService.userService.GetUserCredential(post.AuthorID)
	return blogdto.PostDetailsResponse{
		ID:         post.ID,
		Title:      post.Title,
		Content:    post.Content,
		Status:     uint(post.Status),
		Author:     author.FirstName + " " + author.LastName,
		CoverImage: post.CoverImage,
		CreatedAt:  post.CreatedAt,
	}
}

func (blogService *BlogService) EditPost(request blogdto.EditPostRequest) {
	ok := blogService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		panic(forbiddenError)
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	post, exist := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		panic(notFoundError)
	}

	if request.Title != nil {
		post.Title = *request.Title
	}
	if request.Content != nil {
		post.Content = *request.Content
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

	err := blogService.blogRepository.UpdatePost(blogService.db, post)
	if err != nil {
		panic(err)
	}
}

func (blogService *BlogService) DeletePost(request blogdto.DeletePostRequest) {
	ok := blogService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		panic(forbiddenError)
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	for _, postID := range request.PostIDs {
		_, exist := blogService.blogRepository.FindPostByID(blogService.db, postID)
		if !exist {
			continue
		}
		blogService.blogRepository.DeletePost(blogService.db, postID)
	}
}

func (blogService *BlogService) AddPostMedia(request blogdto.AddPostMediaRequest) uint {
	ok := blogService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		panic(forbiddenError)
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	_, exist := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		panic(notFoundError)
	}

	mediaPath := blogService.constants.S3BucketPath.GetBlogMediaPath(request.PostID, request.Media.Filename)
	blogService.s3Storage.UploadObject(enum.BlogMedia, mediaPath, request.Media)

	media := &entity.Media{
		Path:      mediaPath,
		OwnerID:   request.PostID,
		OwnerType: "blog",
	}
	if err := blogService.blogRepository.AddMedia(blogService.db, media); err != nil {
		panic(err)
	}
	return media.ID
}

func (blogService *BlogService) DeletePostMedia(request blogdto.AccessPostMediaRequest) {
	ok := blogService.userService.IsUserActive(request.AuthorID)
	if !ok {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Post,
		}
		panic(forbiddenError)
	}

	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)

	_, exist := blogService.blogRepository.FindPostByID(blogService.db, request.PostID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Post}
		panic(notFoundError)
	}

	media, exist := blogService.blogRepository.GetMediaByID(blogService.db, request.MediaID)
	if !exist {
		notFoundError := exception.NotFoundError{Item: blogService.constants.Field.Media}
		panic(notFoundError)
	}

	if media.OwnerID != request.PostID {
		forbiddenError := exception.ForbiddenError{
			Message:  "",
			Resource: blogService.constants.Field.Media,
		}
		panic(forbiddenError)
	}

	if err := blogService.s3Storage.DeleteObject(enum.BlogMedia, media.Path); err != nil {
		panic(err)
	}
	if err := blogService.blogRepository.DeleteMedia(blogService.db, request.MediaID); err != nil {
		panic(err)
	}
}
