package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.NotFound = http.HandlerFunc(app.notFoundResponse)
    router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

    router.HandlerFunc(http.MethodPost, "/set", app.createRecordHandler)
    router.HandlerFunc(http.MethodGet, "/get/:key", app.showRecordHandler)

    return app.recoverPanic(router)
}
