package model

import (
	"time"
)

type Post struct {
	ID        int64
	Title     string
	Content   string
	UserID    int64
	Tags      []string
	Comments  []Comment
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostFeed struct {
	Post
	Username      string // author's username
	CommentsCount int
}

// PostPatch holds a partial update. Nil fields are left unchanged.
type PostPatch struct {
	Title   *string
	Content *string
	Tags    *[]string
	// Version is the version the client last read (optimistic locking).
	Version int
}
