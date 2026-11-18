package handler

import (
	"net/http"

	"github.com/adocoder12/social_blog/internal/dto"
)

// GetUserFeedHandler returns the user's own posts plus posts from users they follow.
// The user comes from the URL until auth exists; then it will come from the token.
func (app *Application) GetUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	// Parse errors (limit=abc) are 400; out-of-range values (limit=500) are 422.
	fq, err := dto.ParseFeedQuery(r)
	if err != nil {
		app.badRequest(w, r, err.Error(), err)
		return
	}
	if err := app.validator.Struct(fq); err != nil {
		app.failedValidation(w, r, err)
		return
	}

	feed, err := app.posts.GetUserFeed(r.Context(), user.ID, fq.ToModel())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	response := make([]dto.ResponseFeedPost, 0, len(feed))
	for i := range feed {
		response = append(response, dto.FromFeedModel(&feed[i]))
	}

	app.jsonResponse(w, r, http.StatusOK, response)
}
