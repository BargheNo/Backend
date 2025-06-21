package service

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	CreatePost(request blogdto.CreatePostRequest) (uint, error)
	EditPost(request blogdto.EditPostRequest) error
	GetCorporationPosts(request blogdto.GetCorporationPostsRequest) ([]blogdto.CorporationPostResponse, error)
	GetCorporationPostsForGeneral(request blogdto.GetPublicCorporationPostsRequest) ([]blogdto.GeneralPostResponse, error)
	GetGeneralPosts(request blogdto.GetPublicPostsRequest) ([]blogdto.GeneralPostResponse, error)
	GetCorporationPost(request blogdto.GetCorporationPostRequest) (blogdto.CorporationPostResponse, error)
	GetGeneralPost(postID uint) (blogdto.GeneralPostResponse, error)
	DeletePost(request blogdto.DeletePostRequest) error
	AddPostMedia(request blogdto.AddPostMediaRequest) (uint, error)
	DeletePostMedia(request blogdto.AccessPostMediaRequest) error
	GetPostMedia(request blogdto.AccessPostMediaRequest) (string, error)
	LikePost(request blogdto.LikePostRequest) error
	UnlikePost(request blogdto.LikePostRequest) error
}
