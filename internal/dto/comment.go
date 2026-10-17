package dto

import (
	"time"

	"github.com/adocoder12/social_blog/internal/model"
)

type RequestComment struct {
	Content string `json:"content" validate:"required,max=1000"`
	UserID  int64  `json:"user_id" validate:"required,gt=0"`
}

type ResponseComment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *RequestComment) ToModel(postID int64) *model.Comment {
	return &model.Comment{
		PostID:  postID,
		UserID:  r.UserID,
		Content: r.Content,
	}
}

func CommentFromModel(c *model.Comment) ResponseComment {
	return ResponseComment{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		Username:  c.User.Username,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
	}
}
