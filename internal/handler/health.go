package handler

import "net/http"

func (app *Application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.AppEnv,
		"version": app.version,
	}

	if err := app.writeJSON(w, http.StatusOK, data); err != nil {
		app.serverError(w, r, err)
	}
}
