package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (app *Application) logError(err error) {
	app.Logger.Error("Application error happenned", slog.Any("error", err))
}

// NOTE: it is okay to make separate structure for error response, 
// but I decided not to :)
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request,
	status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(500)
	}
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(err)
	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err *error) {
	if err==nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Bad request was recieved")
	} else {
		app.errorResponse(w, r, http.StatusBadRequest, (*err).Error())
	}
}
