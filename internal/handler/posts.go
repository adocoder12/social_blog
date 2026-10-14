package handler

import (
	"errors"
	"net/http"

	"github.com/adocoder12/social_blog/internal/dto"
	"github.com/adocoder12/social_blog/internal/store"
)

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

	post, err := app.posts.CreatePost(r.Context(), req.Tomodel())
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