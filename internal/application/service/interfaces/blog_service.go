package service

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	CreatePost(request blogdto.CreatePostRequest) uint
	EditPost(request blogdto.EditPostRequest)
	GetCorporationPosts(request blogdto.GetPostsRequest) []blogdto.CorporationPostResponse
	GetCorporationPostsForGeneral(request blogdto.GetPostsRequest) []blogdto.GeneralPostResponse
	GetGeneralPosts(request blogdto.GetPostsRequest) []blogdto.GeneralPostResponse
	GetCorporationPost(request blogdto.GetPostRequest) blogdto.CorporationPostResponse
	GetGeneralPost(request blogdto.GetPostRequest) blogdto.GeneralPostResponse
	DeletePost(request blogdto.DeletePostRequest)
	AddPostMedia(request blogdto.AddPostMediaRequest) uint
	DeletePostMedia(request blogdto.AccessPostMediaRequest)
	GetPostMedia(request blogdto.AccessPostMediaRequest) string
	LikePost(request blogdto.LikePostRequest)
	UnlikePost(request blogdto.LikePostRequest)
}
