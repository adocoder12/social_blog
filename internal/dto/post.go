package dto

import (
	"time"

	"github.com/adocoder12/social_blog/internal/model"
)

type ResponsePost struct {
	ID        int64             `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	UserID    int64             `json:"user_id"`
	Tags      []string          `json:"tags"`
	Comments  []ResponseComment `json:"comments"`
	Version   int				`json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type RequestPost struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	UserID  int64    `json:"user_id" validate:"required,gt=0"`
	Tags    []string `json:"tags" validate:"omitempty,max=10,dive,required,max=30"`
}

type UpdatePost struct {
	Title   *string   `json:"title" validate:"omitempty,min=1,max=100"`
	Content *string   `json:"content" validate:"omitempty,min=1,max=1000"`
	Tags    *[]string `json:"tags" validate:"omitempty,max=10,dive,required,max=30"`
}

func (r *RequestPost) ToModel() *model.Post {
	return &model.Post{
		Title:   r.Title,
		Content: r.Content,
		UserID:  r.UserID,
		Tags:    r.Tags,
	}
}

func (u *UpdatePost) ToPatch() model.PostPatch {
	return model.PostPatch{
		Title:   u.Title,
		Content: u.Content,
		Tags:    u.Tags,
	}
}

func FromModel(p *model.Post) ResponsePost {
	comments := make([]ResponseComment, 0, len(p.Comments))
	for i := range p.Comments {
		comments = append(comments, CommentFromModel(&p.Comments[i]))
	}
	response := ResponsePost{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		UserID:    p.UserID,
		Tags:      p.Tags,
		Comments:  comments,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return response
}
