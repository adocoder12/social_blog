package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/adocoder12/social_blog/internal/dto"
	"github.com/adocoder12/social_blog/internal/model"
	"github.com/adocoder12/social_blog/internal/store"
)

type postKey string

const postCtx postKey = "post"

func (app *Application) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.GetAllPosts(r.Context())

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	response := make([]dto.ResponsePost, 0, len(posts))
	for i := range posts {
		response = append(response, dto.FromModel(&posts[i]))
	}

	if err := app.writeJSON(w, http.StatusOK, response); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RequestPost
	if err := app.readJSON(w, r, &req); err != nil {
		app.badRequest(w, r, "invalid request body", err)
		return
	}

	if err := app.validator.Struct(req); err != nil {
		app.failedValidation(w, r, err)
		return
	}

	post, err := app.posts.CreatePost(r.Context(), req.ToModel())
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			app.badRequest(w, r, "user does not exist", err)
			return
		}
		app.serverError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, dto.FromModel(post)); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.posts.GetComments(r.Context(), post.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.writeJSON(w, http.StatusOK, dto.FromModel(post)); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) PostContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID, err := app.readIDParam(r, "postID")
		if err != nil {
			app.badRequest(w, r, "invalid post id", err)
			return
		}
		post, err := app.posts.GetPostByID(r.Context(), postID)
		if err != nil {
			if errors.Is(err, store.ErrPostNotFound) {
				app.notFound(w, r)
				return
			}
			app.serverError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *model.Post {
	post, _ := r.Context().Value(postCtx).(*model.Post)
	return post
}

func (app *Application) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var req dto.UpdatePost
	if err := app.readJSON(w, r, &req); err != nil {
		app.badRequest(w, r, "invalid request body", err)
		return
	}
	if err := app.validator.Struct(req); err != nil {
		app.failedValidation(w, r, err)
		return
	}

	updated, err := app.posts.UpdatePost(r.Context(), post, req.ToPatch())
	if err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, dto.FromModel(updated)); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	if err := app.posts.DeletePost(r.Context(), post.ID); err != nil {
		if errors.Is(err, store.ErrPostNotFound) {
			app.notFound(w, r)
			return
		}
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
