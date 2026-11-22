package handlers

import (
	"net/http"
)

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.Config.Env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, 200, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
