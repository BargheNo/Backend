package service

import blogdto "github.com/BargheNo/Backend/internal/application/dto/blog"

type BlogService interface {
	CreatePost(request blogdto.CreatePostRequest)
	GetCorporationPosts(request blogdto.GetCorporationPostsRequest) ([]blogdto.PostResponse, error)
	GetPost(postID uint) blogdto.PostDetailsResponse
}
