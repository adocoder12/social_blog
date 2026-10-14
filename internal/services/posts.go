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
