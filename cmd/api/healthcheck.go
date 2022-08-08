package main

import (
	"net/http"
	"strconv"

	"github.com/MinaWilliam/movies/internal/services/visits"
)

func (app *app) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	visits.RecordVisit()
	data := envelope{
		"status": "available",
		"visits": strconv.FormatUint(visits.GetVisitsCount(), 10),
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJson(w, 200, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
