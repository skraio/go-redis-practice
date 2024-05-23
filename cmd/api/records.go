package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/skraio/go-redis-practice/internal/data"
	"github.com/skraio/go-redis-practice/internal/validator"
)

func (app *application) createRecordHandler(w http.ResponseWriter, r *http.Request) {
	var inputRecord data.Record
	err := app.readJSON(w, r, &inputRecord)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateRecord(v, &inputRecord); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.records.Insert(&inputRecord)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "record successfully created"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRecordHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.records.Get(key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"record": record}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRecordHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.records.Delete(key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "record successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) incrRecordValueHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.records.Get(key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.records.Update(record, data.Increase, 1)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "record incremented by 1"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) decrRecordValueHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.records.Get(key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.records.Update(record, data.Decrease, 1)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "record decremented by 1"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) incrbyRecordValueHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.records.Get(key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Term *int64 `json:"term"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Term != nil, "term", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.records.Update(record, data.Increase, *input.Term)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": fmt.Sprintf("record incremented by %d", *input.Term)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) decrbyRecordValueHandler(w http.ResponseWriter, r *http.Request) {
	key, err := app.readKeyParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	record, err := app.records.Get(key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Term *int64 `json:"term"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Term != nil, "term", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.records.Update(record, data.Decrease, *input.Term)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": fmt.Sprintf("record decremented by %d", *input.Term)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
