package service

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	CreatePost(request blogdto.CreatePostRequest)
	EditPost(request blogdto.EditPostRequest)
	GetPosts(request blogdto.GetPostsRequest) []blogdto.PostResponse
	GetPost(request blogdto.GetPostRequest) blogdto.PostDetailsResponse
	DeletePost(request blogdto.DeletePostRequest)
	AddPostMedia(request blogdto.AddPostMediaRequest) uint
	DeletePostMedia(request blogdto.AccessPostMediaRequest)
}
