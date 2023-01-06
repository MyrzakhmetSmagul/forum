package service

import (
	"log"

	"github.com/MyrzakhmetSmagul/forum/internal/model"
	"github.com/MyrzakhmetSmagul/forum/internal/repository"
)

type PostService interface {
	CreatePost(post *model.Post) error
	GetPost(post *model.Post) error
	CreateComment(comment *model.Comment) error
	GetPostComments(post *model.Post) error
	GetAllPosts() ([]model.Post, error)
	PostLike(reaction *model.PostReaction) error
	PostDislike(reaction *model.PostReaction) error
	CommenSettLike(reaction *model.CommentReaction) error
	CommentSetDislike(reaction *model.CommentReaction) error
	GetAllCategory() ([]model.Category, error)
}

type postService struct {
	repository.PostQuery
	repository.PostReactionQuery
	repository.CommentQuery
	repository.CommentReactionQuery
	repository.CategoryQuery
}

func NewPostService(dao repository.DAO) PostService {
	return &postService{
		PostQuery:            dao.NewPostQuery(),
		PostReactionQuery:    dao.NewPostReactionQuery(),
		CommentQuery:         dao.NewCommentQuery(),
		CommentReactionQuery: dao.NewCommentReactionQuery(),
		CategoryQuery:        dao.NewCategoryQuery(),
	}
}

func (p *postService) CreatePost(post *model.Post) error {
	return p.PostQuery.CreatePost(post)
}

func (p *postService) GetPost(post *model.Post) error {
	err := p.PostQuery.GetPost(post)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = p.PostReactionQuery.GetPostReactions(post)
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (p *postService) CreateComment(comment *model.Comment) error {
	return p.CommentQuery.CreateComment(comment)
}

func (p *postService) GetPostComments(post *model.Post) error {
	err := p.CommentQuery.GetPostComments(post)
	if err != nil {
		log.Println(err)
		return err
	}

	for i := 0; i < len(post.Comments); i++ {
		err = p.CommentReactionQuery.GetCommentLikesDislikes(&post.Comments[i])
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (p *postService) GetAllPosts() ([]model.Post, error) {
	return p.PostQuery.GetAllPosts()
}

func (p *postService) PostLike(reaction *model.PostReaction) error {
	return p.PostReactionQuery.PostSetLike(reaction)
}

func (p *postService) PostDislike(reaction *model.PostReaction) error {
	return p.PostReactionQuery.PostSetDislike(reaction)
}

func (p *postService) CommenSettLike(reaction *model.CommentReaction) error {
	return p.CommentReactionQuery.CommentSetLike(reaction)
}

func (p *postService) CommentSetDislike(reaction *model.CommentReaction) error {
	return p.CommentReactionQuery.CommentSetDislike(reaction)
}

func (p *postService) GetAllCategory() ([]model.Category, error) {
	return p.CategoryQuery.GetAllCategory()
}
