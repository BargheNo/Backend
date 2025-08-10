package usecase

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	GetBlogSortableColumns() []blogdto.GetBlogEnumResponse
	CreatePost(request blogdto.CreatePostRequest) (uint, error)
	EditPost(request blogdto.EditPostRequest) error
	GetCorporationPosts(request blogdto.GetCorporationPostsRequest) ([]blogdto.CorporationPostResponse, int64, error)
	GetCorporationPostsForGeneral(request blogdto.GetPublicCorporationPostsRequest) ([]blogdto.GeneralPostResponse, int64, error)
	GetGeneralPosts(request blogdto.GetPublicPostsRequest) ([]blogdto.GeneralPostResponse, int64, error)
	GetCorporationPost(request blogdto.GetCorporationPostRequest) (blogdto.CorporationPostResponse, error)
	GetGeneralPost(postID uint) (blogdto.GeneralPostResponse, error)
	DeletePost(request blogdto.DeletePostRequest) error
	AddPostMedia(request blogdto.AddPostMediaRequest) (uint, error)
	DeletePostMedia(request blogdto.AccessPostMediaRequest) error
	GetPostMedia(request blogdto.AccessPostMediaRequest) (string, error)
	LikePost(request blogdto.GetPostRequest) error
	UnlikePost(request blogdto.GetPostRequest) error
	IsBlogLiked(request blogdto.GetPostRequest) (bool, error)
}
