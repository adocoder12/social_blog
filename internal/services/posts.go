package services

import (
	"context"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/adocoder12/social_blog/internal/store"
)

type PostService struct {
	repo store.PostRepository
}

func NewPostService(repo store.PostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	posts, err := s.repo.GetAllPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("post service: get all posts: %w", err)
	}
	return posts, nil
}

func (s *PostService) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	created, err := s.repo.CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("post service: create post: %w", err)
	}
	return created, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*model.Post, error){
	post,err:= s.repo.GetPostByID(ctx,id)
	if err != nil{
		return  nil, fmt.Errorf("post service: get by id: %w",err)
	}
	return post , nil
}


 
func (s *PostService) UpdatePost(ctx context.Context, id int64, patch model.PostPatch) (*model.Post, error) {
	post, err := s.repo.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("post service: update post: %w", err)
	}

	if patch.Title != nil {
		post.Title = *patch.Title
	}
	if patch.Content != nil {
		post.Content = *patch.Content
	}
	if patch.Tags != nil {
		post.Tags = *patch.Tags
	}

	updated, err := s.repo.UpdatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("post service: update post: %w", err)
	}
	return updated, nil
}
	
func (s *PostService)DeletePost(ctx context.Context, id int64) error{
	if err:= s.repo.DeletePost(ctx,id);err != nil{
		return  fmt.Errorf("post service: delete post: %w",err)
	}
	return nil
}