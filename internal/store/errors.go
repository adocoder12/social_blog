package store

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

// Postgres error codes.
const (
	pgUniqueViolation     = "23505"
	pgForeignKeyViolation = "23503"
	QueryTimeoutDuration = time.Second * 5
)

// isUniqueViolation reports whether err is a unique constraint violation,
// and returns the constraint name.
func isUniqueViolation(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
		return pgErr.ConstraintName, true
	}
	return "", false
}

// isForeignKeyViolation reports whether err is a foreign key violation,
// and returns the constraint name.
func isForeignKeyViolation(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation {
		return pgErr.ConstraintName, true
	}
	return "", false
}
