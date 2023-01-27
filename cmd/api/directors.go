package main

import (
	"fmt"
	"greenlight.aslan/internal/data"
	"greenlight.aslan/internal/validator"
	"net/http"
)

func (app *application) createDirectorHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string   `json:"name"`
		Surname string   `json:"surname"`
		Awards  []string `json:"awards"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	director := &data.Director{
		Name:    input.Name,
		Surname: input.Surname,
		Awards:  input.Awards,
	}

	err = app.models.Directors.Insert(director)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/directors/%d", director.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"director": director}, headers)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDirectorsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string
		Surname string
		Awards  []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Awards = app.readCSV(qs, "awards", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "awards", "fullname",
		"-id", "-awards", "-fullname"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	directors, metadata, err := app.models.Directors.GetAll(input.Name, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"directors": directors, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
