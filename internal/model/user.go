package model

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Hash []byte
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

// Set hashes the plain-text password.
func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.Hash = hash
	return nil
}

// Matches reports whether text matches the stored hash.
func (p *Password) Matches(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(text))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}