package model

import "time"

// FeedQuery holds pagination and filter options for the feed.
// Plain domain type: no HTTP, no JSON, no validation tags.
type FeedQuery struct {
	Limit  int
	Offset int
	Sort   string // "asc" or "desc"
	Tags   []string
	Search string
	Since  *time.Time // nil = no lower bound
	Until  *time.Time // nil = no upper bound
}
