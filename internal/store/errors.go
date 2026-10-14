package store

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (

	// Duplicate unique error
	pgUniqueViolation = "23505"
	// Postgres error code for foreign key violation.
	pgForeignKeyViolation = "23503"
)


func isUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation
}


func isForeignKeyViolationError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation
}