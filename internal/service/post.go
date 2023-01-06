package service

import "github.com/MyrzakhmetSmagul/forum/internal/repository"

type PostService interface {
}

type postService struct {
	repository.PostQuery
	repository.PostReactionQuery
}
