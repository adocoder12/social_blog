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
