package services

import (
	"context"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/adocoder12/social_blog/internal/store"
)

type PostService struct {
	posts    store.PostRepository
	comments store.CommentRepository
}

func NewPostService(posts store.PostRepository, comments store.CommentRepository) *PostService {
	return &PostService{
		posts:    posts,
		comments: comments,
	}
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	posts, err := s.posts.GetAllPosts(ctx)
	if err != nil {
		return nil, fmt.Errorf("post service: get all posts: %w", err)
	}
	return posts, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	post, err := s.posts.GetPostByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("post service: get by id: %w", err)
	}
	return post, nil
}

func (s *PostService) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	created, err := s.posts.CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("post service: create post: %w", err)
	}
	return created, nil
}

func (s *PostService) GetComments(ctx context.Context, postID int64) ([]model.Comment, error) {
	comments, err := s.comments.GetByPostID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("post service: get comments: %w", err)
	}
	return comments, nil
}

func (s *PostService) UpdatePost(ctx context.Context, post *model.Post, patch model.PostPatch) (*model.Post, error) {
	if patch.Title != nil {
		post.Title = *patch.Title
	}
	if patch.Content != nil {
		post.Content = *patch.Content
	}
	if patch.Tags != nil {
		post.Tags = *patch.Tags
	}

	updated, err := s.posts.UpdatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("post service: update post: %w", err)
	}
	return updated, nil
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	if err := s.posts.DeletePost(ctx, id); err != nil {
		return fmt.Errorf("post service: delete post: %w", err)
	}
	return nil
}
