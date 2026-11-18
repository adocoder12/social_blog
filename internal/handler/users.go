package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/adocoder12/social_blog/internal/dto"
	"github.com/adocoder12/social_blog/internal/model"
	"github.com/adocoder12/social_blog/internal/services"
	"github.com/adocoder12/social_blog/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type userKey string

const userCtx userKey = "user"

func (app *Application) GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	app.jsonResponse(w, r, http.StatusOK, dto.FromUserModel(user))
}

func (app *Application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload dto.RequestUser
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, "invalid request body", err)
		return
	}

	if err := app.validator.Struct(payload); err != nil {
		app.failedValidation(w, r, err)
		return
	}
	userModel, err := payload.ToModel()
	if err != nil {
		// validate:"max=72" counts characters; bcrypt's limit is 72 bytes.
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			app.badRequest(w, r, "password must be at most 72 bytes", err)
			return
		}
		app.serverError(w, r, err)
		return
	}
	user, err := app.users.CreateUser(r.Context(), userModel)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.conflict(w, r, "email already taken")
		case errors.Is(err, store.ErrDuplicateUsername):
			app.conflict(w, r, "username already taken")
		default:
			app.serverError(w, r, err)
		}
		return
	}

	app.jsonResponse(w, r, http.StatusCreated, dto.FromUserModel(user))

}

func (app *Application) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var payload dto.UpdateUser
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, "invalid request body", err)
		return
	}
	if err := app.validator.Struct(payload); err != nil {
		app.failedValidation(w, r, err)
		return
	}
	updated, err := app.users.UpdateUser(r.Context(), user, payload.ToPatch())
	if err != nil {
		switch {
		case errors.Is(err, store.ErrUserNotFound):
			app.notFound(w, r)
		case errors.Is(err, store.ErrDuplicateEmail):
			app.conflict(w, r, "email already taken")
		case errors.Is(err, store.ErrDuplicateUsername):
			app.conflict(w, r, "username already taken")
		default:
			app.serverError(w, r, err)
		}
		return
	}
	app.jsonResponse(w, r, http.StatusOK, dto.FromUserModel(updated))

}

func (app *Application) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.users.DeleteUser(r.Context(), user.ID); err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// FollowUserHandler makes the user in the body follow the user in the URL.
func (app *Application) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	target := getUserFromCtx(r)

	var payload dto.FollowUser
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, "invalid request body", err)
		return
	}
	if err := app.validator.Struct(payload); err != nil {
		app.failedValidation(w, r, err)
		return
	}

	if err := app.users.Follow(r.Context(), payload.UserID, target.ID); err != nil {
		switch {
		case errors.Is(err, services.ErrSelfFollow):
			app.badRequest(w, r, "you cannot follow yourself", err)
		case errors.Is(err, store.ErrAlreadyFollowing):
			app.conflict(w, r, "already following this user")
		case errors.Is(err, store.ErrUserNotFound):
			app.badRequest(w, r, "user does not exist", err)
		default:
			app.serverError(w, r, err)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUserHandler makes the user in the body unfollow the user in the URL.
func (app *Application) UnfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	target := getUserFromCtx(r)

	var payload dto.FollowUser
	if err := app.readJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, "invalid request body", err)
		return
	}
	if err := app.validator.Struct(payload); err != nil {
		app.failedValidation(w, r, err)
		return
	}

	if err := app.users.Unfollow(r.Context(), payload.UserID, target.ID); err != nil {
		if errors.Is(err, store.ErrNotFollowing) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// middlewares
func (app *Application) UserContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := app.readIDParam(r, "userID")
			if err != nil {
				app.badRequest(w, r, "invalid user id", err)
				return
			}

			user, err := app.users.GetUserByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, store.ErrUserNotFound) {
					app.notFound(w, r)
					return
				}
				app.serverError(w, r, err)
				return
			}
			ctx := context.WithValue(r.Context(), userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})

}

func getUserFromCtx(r *http.Request) *model.User {
	// Unchecked on purpose: a missing value means the middleware isn't wired,
	// a programmer error. Recoverer turns the panic into a 500.
	return r.Context().Value(userCtx).(*model.User)
}
