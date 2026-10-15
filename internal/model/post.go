package model

import "time"

type Post struct {
	ID        int64
	Title     string
	Content   string
	UserID    int64
	Tags      []string
	Comments  []Comment
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PostPatch holds a partial update. Nil fields are left unchanged.
type PostPatch struct {
	Title   *string
	Content *string
	Tags    *[]string
}
