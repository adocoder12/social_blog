package dto

import (
	"time"

	"github.com/adocoder12/social_blog/internal/model"
)

type ResponseUser struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

type RequestUser struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type UpdateUser struct {
	Username *string `json:"username" validate:"omitempty,min=1,max=100"`
	Email    *string `json:"email" validate:"omitempty,email,max=255"`
}

func (r *RequestUser) ToModel() (*model.User, error) {
	user := &model.User{
		Username: r.Username,
		Email:    r.Email,
	}
	if err := user.Password.Set(r.Password); err != nil {
		return nil, err
	}
	return user, nil
}

func FromUserModel(u *model.User) ResponseUser {
	return ResponseUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}
}

func (u *UpdateUser) ToPatch() model.UserPatch {
	return model.UserPatch{
		Username: u.Username,
		Email:    u.Email,
	}
}

// FollowUser identifies who is following. Temporary: once auth exists,
// the follower comes from the authenticated user, not the body.
type FollowUser struct {
	UserID int64 `json:"user_id" validate:"required,gt=0"`
}
