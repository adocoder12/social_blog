package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrEditConflict = errors.New("edit conflict: resource was modified")
)

type PostStore struct {
	db DBTX
}

// sortDirections maps the allowed sort values to SQL. ORDER BY direction can't be
// a $N placeholder, so it goes into the query string — only via this whitelist.
var sortDirections = map[string]string{
	"asc":  "ASC",
	"desc": "DESC",
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq model.FeedQuery) ([]model.PostFeed, error) {
	dir, ok := sortDirections[fq.Sort]
	if !ok {
		dir = "DESC"
	}

	// nil would be sent as NULL, and "tags @> NULL" filters out every row.
	tags := fq.Tags
	if tags == nil {
		tags = []string{}
	}

	// Own posts + posts from users that userID follows, then optional filters.
	// Each filter is "param is empty OR condition", so unset filters match everything.
	// p.id breaks ties so rows don't repeat or vanish between pages.
	//
	// followers: user_id = the followed user, follower_id = who follows them.
	// The follow check is a subquery, not a JOIN: a JOIN would repeat each post
	// once per matching follower row and inflate COUNT(c.id).
	// users is an INNER JOIN: posts.user_id is a NOT NULL FK, so the author always exists.
	// GROUP BY p.id, u.id is enough: other columns depend on those primary keys.
	query := `
		SELECT
			p.id, p.title, p.content, p.user_id, p.tags, p.version, p.created_at, p.updated_at,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		JOIN users u ON u.id = p.user_id
		LEFT JOIN comments c ON c.post_id = p.id
		WHERE (p.user_id = $1
		       OR p.user_id IN (SELECT user_id FROM followers WHERE follower_id = $1))
		  AND (cardinality($4::text[]) = 0 OR p.tags @> $4)
		  AND ($5 = '' OR p.title ILIKE '%' || $5 || '%' OR p.content ILIKE '%' || $5 || '%')
		  AND ($6::timestamptz IS NULL OR p.created_at >= $6)
		  AND ($7::timestamptz IS NULL OR p.created_at <= $7)
		GROUP BY p.id, u.id
		ORDER BY p.created_at ` + dir + `, p.id ` + dir + `
		LIMIT $2 OFFSET $3`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	// Argument order must match $1..$7 above.
	rows, err := s.db.Query(ctx, query,
		userID,    // $1
		fq.Limit,  // $2
		fq.Offset, // $3
		tags,      // $4
		fq.Search, // $5
		fq.Since,  // $6
		fq.Until,  // $7
	)
	if err != nil {
		return nil, fmt.Errorf("post repo - get user feed query: %w", err)
	}
	defer rows.Close()

	feed := []model.PostFeed{}
	for rows.Next() {
		var post model.PostFeed
		// Scan order must match the SELECT column order.
		if err := rows.Scan(
			&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags, &post.Version, &post.CreatedAt, &post.UpdatedAt,
			&post.Username,
			&post.CommentsCount,
		); err != nil {
			return nil, fmt.Errorf("post repo - get user feed scan: %w", err)
		}
		feed = append(feed, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("post repo - get user feed rows: %w", err)
	}
	return feed, nil
}

func (s *PostStore) GetAllPosts(ctx context.Context) ([]model.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, version ,created_at, updated_at
		FROM posts
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("post repo - get all query: %w", err)
	}
	defer rows.Close()

	posts := []model.Post{}
	for rows.Next() {
		var post model.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags, &post.Version, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return nil, fmt.Errorf("post repo - get all scan: %w", err)
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("post repo - get all rows: %w", err)
	}
	return posts, nil
}

func (s *PostStore) GetPostByID(ctx context.Context, id int64) (*model.Post, error) {
	query := `
		SELECT id, title, content, user_id, tags, version, created_at, updated_at
		FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post model.Post
	err := s.db.QueryRow(ctx, query, id).
		Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.Tags, &post.Version, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, fmt.Errorf("post repo - get by id query row: %w", err)
	}
	return &post, nil
}

func (s *PostStore) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	query := `
		INSERT INTO posts (title, content, user_id, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`
	// pgx v5 handles Go slices natively; nil would be sent as NULL and violate NOT NULL.
	if post.Tags == nil {
		post.Tags = []string{}
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, post.Title, post.Content, post.UserID, post.Tags).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if _, ok := isForeignKeyViolation(err); ok {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("post repo - create post query: %w", err)
	}
	return post, nil
}

func (s *PostStore) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, updated_at = now() , version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING user_id, version, updated_at`

	if post.Tags == nil {
		post.Tags = []string{}
	}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRow(ctx, query, post.Title, post.Content, post.Tags, post.ID, post.Version).
		Scan(&post.UserID, &post.Version, &post.UpdatedAt)
	if err != nil {
		// No row matched id+version: the post changed (or was deleted) since the client read it.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEditConflict
		}
		return nil, fmt.Errorf("post repo - update query: %w", err)
	}
	return post, nil
}

func (s *PostStore) DeletePost(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("post repo - delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrPostNotFound
	}
	return nil
}
