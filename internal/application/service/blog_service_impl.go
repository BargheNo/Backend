package serviceimpl

import (
	"github.com/BargheNo/Backend/bootstrap"
	blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"
	service "github.com/BargheNo/Backend/internal/application/service/interfaces"
	"github.com/BargheNo/Backend/internal/domain/entity"
	"github.com/BargheNo/Backend/internal/domain/enum"
	repository "github.com/BargheNo/Backend/internal/domain/repository/postgres"
	"github.com/BargheNo/Backend/internal/domain/s3"
	"github.com/BargheNo/Backend/internal/infrastructure/database"
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
	blogService.corporationService.CheckApplicantAccess(request.CorporationID, request.AuthorID)
	println(request.CoverImage.Filename)
	post := &entity.Post{
		Title:         request.Title,
		Content:       request.Content,
		AuthorID:      request.AuthorID,
		CorporationID: request.CorporationID,
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
