package main

import (
	"fmt"
	"net/http"

	"github.com/skraio/mini-godis/internal/data"
	"github.com/skraio/mini-godis/internal/validator"
)

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Key   string `json:"key"`
		Value int32 `json:"value"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	record := &data.Record{
		Key:   input.Key,
		Value: input.Value,
	}

	v := validator.New()
	if data.ValidateRecord(v, record); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record := data.Record{
		Key:   key,
		Value: 16,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
