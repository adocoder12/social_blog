package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

// errorResponse sends {"error": msg} with the given status.
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, msg string) {
	if err := app.writeJSON(w, status, map[string]string{"error": msg}); err != nil {
		app.logger.Error("writte error response",
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