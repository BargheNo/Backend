package service

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	CreatePost(request blogdto.CreatePostRequest) error
	EditPost(request blogdto.EditPostRequest) error
	GetCorporationPosts(request blogdto.GetPostsRequest) ([]blogdto.CorporationPostResponse, error)
	GetCorporationPostsForGeneral(request blogdto.GetPostsRequest) ([]blogdto.GeneralPostResponse, error)
	GetGeneralPosts(request blogdto.GetPostsRequest) ([]blogdto.GeneralPostResponse, error)
	GetCorporationPost(request blogdto.GetPostRequest) (blogdto.CorporationPostResponse, error)
	GetGeneralPost(request blogdto.GetPostRequest) (blogdto.GeneralPostResponse, error)
	DeletePost(request blogdto.DeletePostRequest) error
	AddPostMedia(request blogdto.AddPostMediaRequest) (uint, error)
	DeletePostMedia(request blogdto.AccessPostMediaRequest) error
	GetPostMedia(request blogdto.AccessPostMediaRequest) (string, error)
	LikePost(request blogdto.LikePostRequest) error
	UnlikePost(request blogdto.LikePostRequest) error
}
