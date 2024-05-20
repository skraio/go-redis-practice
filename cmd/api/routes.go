package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/record", app.createRecordHandler)
	router.HandlerFunc(http.MethodGet, "/record/:key", app.showRecordHandler)
	router.HandlerFunc(http.MethodDelete, "/record/:key", app.deleteRecordHandler)

	router.HandlerFunc(http.MethodPut, "/record/:key/increment", app.incrRecordValueHandler)
	router.HandlerFunc(http.MethodPut, "/record/:key/decrement", app.decrRecordValueHandler)
	router.HandlerFunc(http.MethodPut, "/record/:key/increment-by", app.incrbyRecordValueHandler)
	router.HandlerFunc(http.MethodPut, "/record/:key/decrement-by", app.decrbyRecordValueHandler)

	return app.recoverPanic(router)
}
