package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/adocoder12/social_blog/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pgUniqueViolation = "23505"

var (
	ErrUserExists   = errors.New("email or username already exists")
	ErrUserNotFound = errors.New("user not found")
)


// helper 
func isUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}


type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}



func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query, user.Email, user.Username, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("user repo - create: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	query := `SELECT id, email, username, created_at, updated_at FROM users`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("user repo - get all query: %w", err)
	}
	defer rows.Close()

	var users []model.User
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

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	query := `SELECT id, email, username, created_at, updated_at FROM users WHERE id = $1`

	var user model.User
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("user repo - get by id: %w", err)
	}
	return &user, nil
}
func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET username = $1,
		    email = $2,
		    updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query, user.Username, user.Email, user.ID).
		Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		if isUniqueViolationError(err) {
			return nil, ErrUserExists
		}
		return nil, fmt.Errorf("user repo - update: %w", err)
	}
	return user, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int, hash string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.pool.Exec(ctx, query, hash, id)
	if err != nil {
		return fmt.Errorf("user repo - update password: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("user repo - delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}