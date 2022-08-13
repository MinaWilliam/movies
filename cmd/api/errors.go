package main

import (
	"fmt"
	"net/http"
)

func (app *app) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

func (app *app) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {

	err := app.writeJson(w, status, envelope{"error": message}, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *app) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, r, http.StatusInternalServerError, "The server encountered a problem and could not process your request")
}

func (app *app) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *app) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *app) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusConflict, "unable to update the record due to an edit conflict, please try again")
}

func (app *app) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, "The requested resource could not be found")
}

func (app *app) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("the %s method is not supported for this resource", r.Method))
}
