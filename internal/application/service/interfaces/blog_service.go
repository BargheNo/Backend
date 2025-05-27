package service

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	CreatePost(request blogdto.CreatePostRequest)
	EditPost(request blogdto.EditPostRequest)
	GetCorporationPosts(request blogdto.GetCorporationPostsRequest) ([]blogdto.PostResponse, error)
	GetPost(postID uint) blogdto.PostDetailsResponse
	DeletePost(request blogdto.DeletePostRequest)
	AddPostMedia(request blogdto.AddPostMediaRequest) uint
	DeletePostMedia(request blogdto.AccessPostMediaRequest)
}
