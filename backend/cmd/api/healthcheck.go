package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	healthcheckData := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version": version,
		},
	}

	err := app.writeJson(w, http.StatusOK, healthcheckData, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
