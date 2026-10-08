package model

import "time"

type User struct{
	ID int64
	Username string
	Email string
	Password []byte
	CreatedAt time.Time
	UpdatedAt time.Time

}

