package dto

import (
	"time"

	"github.com/adocoder12/social_blog/internal/model"
)

type ResponsePost struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestPost struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	UserID  int64     `json:"user_id"`
	Tags    *[]string `json:"tags"`
}

type UpdatePost struct {
	Title   *string   `json:"title"`
	Content *string   `json:"content"`
	Tags    *[]string `json:"tags"`
}

func (r *RequestPost) Tomodel() *model.Post {
	return &model.Post{
		Title:   r.Title,
		Content: r.Content,
		UserID:  r.UserID,
		Tags:    *r.Tags,
	}
}

func fromModel(p *model.Post) ResponsePost {
	response := ResponsePost{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		UserID:    p.UserID,
		Tags:      p.Tags,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
	return response
}
