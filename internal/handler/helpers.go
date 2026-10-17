package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func (app *Application) writeJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("encode response: %w", err)
	}
	body = append(body, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		app.logger.Warn("write response", "err", err)
	}
	return nil
}

func (app *Application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1_048_576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	dec := json.NewDecoder(r.Body)
	// Rejects fields the struct doesn't have, such as typos from the client.
	dec.DisallowUnknownFields()

	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("decode request body: %w", err)
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return errors.New("body must contain a single JSON value")
	}

	return nil
}

// readIDParam parses a positive int64 ID from the URL path.
func (app *Application) readIDParam(r *http.Request, paramName string) (int64, error) {
	strID := chi.URLParam(r, paramName)
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil || id < 1 {
		return 0, fmt.Errorf("invalid %s: %q", paramName, strID)
	}
	return id, nil
}

// errorResponse sends {"error": msg} with the given status.
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, msg string) {
	if err := app.writeJSON(w, status, map[string]string{"error": msg}); err != nil {
		app.logger.Error("write error response",
			"err", err,
			"method", r.Method,
			"uri", r.URL.RequestURI(),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverError logs the real error and sends a generic 500 to the client.
func (app *Application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("server error",
		"err", err,
		"method", r.Method,
		"uri", r.URL.RequestURI(),
	)
	app.errorResponse(w, r, http.StatusInternalServerError, "the server encountered a problem")
}

// notFound sends a 404.
func (app *Application) notFound(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("not found",
		"method", r.Method,
		"uri", r.URL.RequestURI(),
	)
	app.errorResponse(w, r, http.StatusNotFound, "not found")
}

// methodNotAllowed sends a 405.
func (app *Application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, "method not allowed")
}

// badRequest logs the problem and sends a 400 with the error message.
// Always pass a non-nil error.
func (app *Application) badRequest(w http.ResponseWriter, r *http.Request, msg string, err error) {
	app.logger.Warn("bad request",
		"msg", msg,
		"err", err,
		"method", r.Method,
		"uri", r.URL.RequestURI(),
	)
	app.errorResponse(w, r, http.StatusBadRequest, msg)
}

// failedValidation sends a 422 with {"errors": {"field": "rule"}}.
func (app *Application) failedValidation(w http.ResponseWriter, r *http.Request, err error) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		app.badRequest(w, r, "invalid request body", err)
		return
	}

	fields := make(map[string]string, len(ve))
	for _, fe := range ve {
		if fe.Param() != "" {
			fields[fe.Field()] = fe.Tag() + "=" + fe.Param()
		} else {
			fields[fe.Field()] = fe.Tag()
		}
	}

	if err := app.writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": fields}); err != nil {
		app.serverError(w, r, err)
	}
}
