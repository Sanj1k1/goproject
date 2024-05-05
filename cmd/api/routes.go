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
	
	router.HandlerFunc(http.MethodGet, "/v1/characters", app.requirePermission("characters:read",app.listCharactersHandler))
	router.HandlerFunc(http.MethodPost, "/v1/characters", app.requirePermission("characters:write",app.createCharacterHandler))
	router.HandlerFunc(http.MethodGet, "/v1/characters/:id", app.requirePermission("characters:read",app.showCharacterHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id", app.requirePermission("characters:write",app.updateCharacterHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/characters/:id", app.requirePermission("characters:write",app.deleteCharacterHandler))
	
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
