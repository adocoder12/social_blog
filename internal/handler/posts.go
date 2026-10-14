package handler

import (
	"net/http"

	"github.com/adocoder12/social_blog/internal/dto"
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
