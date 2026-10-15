package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserExists        = errors.New("email or username already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already in use")
	ErrDuplicateUsername = errors.New("username already taken")
)

type UsersStore struct {
	pool *pgxpool.Pool
}

func NewUsersStore(pool *pgxpool.Pool) *UsersStore {
	return &UsersStore{pool: pool}
}

// userUniqueError maps a unique violation on users to a specific sentinel error.
func userUniqueError(constraint string) error {
	switch constraint {
	case "users_email_key":
		return ErrDuplicateEmail
	case "users_username_key":
		return ErrDuplicateUsername
	}
	return ErrUserExists
}

func (s *UsersStore) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, user.Email, user.Username, user.Password.Hash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if constraint, ok := isUniqueViolation(err); ok {
			return nil, userUniqueError(constraint)
		}
		return nil, fmt.Errorf("user repo - create: %w", err)
	}
	return user, nil
}

func (s *UsersStore) GetAllUsers(ctx context.Context) ([]model.User, error) {
	query := `
		SELECT id, email, username, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("user repo - get all query: %w", err)
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Username, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("user repo - get all scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user repo - get all rows: %w", err)
	}

	return users, nil
}

func (s *UsersStore) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, username, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user model.User
	err := s.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("user repo - get by id: %w", err)
	}
	return &user, nil
}

func (s *UsersStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, username, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user model.User
	err := s.pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("user repo - get by email: %w", err)
	}
	return &user, nil
}

func (s *UsersStore) UpdateUser(ctx context.Context, user *model.User) (*model.User, error){
	query := `
		UPDATE users
		SET username = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	err := s.pool.QueryRow(ctx, query, user.Username, user.Email, user.ID).
		Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		if constraint, ok := isUniqueViolation(err); ok {
			return nil, userUniqueError(constraint)
		}
		return nil, fmt.Errorf("user repo - update: %w", err)
	}
	return user, nil
}

func (s *UsersStore) UpdatePassword(ctx context.Context, id int64, hash []byte) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`

	result, err := s.pool.Exec(ctx, query, hash, id)
	if err != nil {
		return fmt.Errorf("user repo - update password: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *UsersStore) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("user repo - delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}