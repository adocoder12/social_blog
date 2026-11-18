package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash []byte
}

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}

// UserPatch holds a partial update. Nil fields are left unchanged.
type UserPatch struct {
	Username *string
	Email    *string
}

// Set hashes the plain-text password.
func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.hash = hash
	return nil
}

// Hash returns the bcrypt hash for storage.
func (p *Password) Hash() []byte {
	return p.hash
}

// SetHash loads an existing bcrypt hash, e.g. one read from the database.
func (p *Password) SetHash(hash []byte) {
	p.hash = hash
}

// Matches reports whether text matches the stored hash.
func (p *Password) Matches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(text))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
