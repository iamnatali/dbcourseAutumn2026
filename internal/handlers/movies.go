package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"hotel.booking.natali/internal/domain"
)

// todo:
// сейчас createMovieHandler зависит от Repository напрямую
// а должен зависеть от MoviewService,
// в котором будет объявлена своя movieRepo с урезанным функционалом
func (app *Application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string         `json:"title"`
		ReleaseDate domain.Date    `json:"releaseDate"`
		Runtime     domain.Runtime `json:"runtime"`
		Genres      []string       `json:"genres"`
		Language    string         `json:"language"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, &err)
		return
	}

	movie := &domain.Movie{
		Title:       input.Title,
		ReleaseDate: input.ReleaseDate,
		Runtime:     input.Runtime,
		Genres:      input.Genres,
		Language:    input.Language,
	}

	err = (*app.Repository).InsertMovie(r.Context(), movie)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie},
		headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	successFlag, id := app.parseId(w, r)
	if !successFlag {
		return
	}

	movie, err := (*app.Repository).DeleteMovie(r.Context(), id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	} else {
		err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}

}

func (app *Application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	//todo: cut into separate funtion
	var input struct {
		Title       string         `json:"title"`
		ReleaseDate domain.Date    `json:"releaseDate"`
		Runtime     domain.Runtime `json:"runtime"`
		Genres      []string       `json:"genres"`
		Language    string         `json:"language"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, &err)
		return
	}

	successFlag, id := app.parseId(w, r)
	if !successFlag {
		return
	}

	movie := domain.Movie{
		ID:          id,
		Title:       input.Title,
		ReleaseDate: input.ReleaseDate,
		Runtime:     input.Runtime,
		Genres:      input.Genres,
		Language:    input.Language,
	}

	//todo: возвращать в методах repository не ошибки pgx напрямую, 
	//а оборачивать их в структуры кастомных ошибок 
	//(структуры с методами Error & Unwrap)
	//https://gist.github.com/renanbastos93/f63538b4e75d0bd0fc2d56537206f186
	//https://go.dev/blog/go1.13-errors 
	err = (*(app.Repository)).UpdateMovie(r.Context(), &movie)

	if err != nil {
		//todo: в случае, если ошибка заключается в том, 
		//что фильм не найден, нужно возвращать пользователю не 500, а 404
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) parseId(w http.ResponseWriter, r *http.Request) (bool, uuid.UUID) {
	strID := chi.URLParam(r, "id")
	id, err := uuid.Parse(strID)
	if err != nil {
		app.badRequestResponse(w, r, nil)
		return false, uuid.Nil
	} else {
		return true, id
	}
}

func (app *Application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	successFlag, id := app.parseId(w, r)
	if !successFlag {
		return
	}

	movies, err := (*app.Repository).GetMovie(r.Context(), id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	} else {
		err = app.writeJSON(w, http.StatusOK, envelope{"movie": movies}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
	}
}
