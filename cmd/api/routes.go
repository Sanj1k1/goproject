package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	
	router.HandlerFunc(http.MethodGet, "/v1/characters", app.listCharactersHandler)
	router.HandlerFunc(http.MethodPost, "/v1/characters", app.createCharacterHandler)
	router.HandlerFunc(http.MethodGet, "/v1/characters/:id", app.showCharacterHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id", app.updateCharacterHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/characters/:id", app.deleteCharacterHandler)
	
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
