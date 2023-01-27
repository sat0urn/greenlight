package main

import (
	"fmt"
	"greenlight.aslan/internal/data"
	"net/http"
)

func (app *application) createTrailerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Trailer_name string `json:"trailer_name"`
		Duration     int64  `json:"duration"`
		Premier_date string `json:"premier_date"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	trailer := &data.Trailer{
		Trailer_name: input.Trailer_name,
		Duration:     input.Duration,
		Premier_date: input.Premier_date,
	}

	err = app.models.Trailers.Insert(trailer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/trailers/%d", trailer.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"trailer": trailer}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
