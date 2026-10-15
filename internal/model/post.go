package model

import "time"

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostPatch holds a partial update. Nil fields are left unchanged.
type PostPatch struct {
	Title   *string
	Content *string
	Tags    *[]string
}