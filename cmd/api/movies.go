package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MinaWilliam/movies/internal/data"
)

func (app *app) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "list movies")
}

func (app *app) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    uint32       `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *app) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.getIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		Title:     "Casablanca",
		Runtime:   102,
		Year:      1996,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
		CreatedAt: time.Now(),
	}

	err = app.writeJson(w, http.StatusOK, envelope{"data": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
