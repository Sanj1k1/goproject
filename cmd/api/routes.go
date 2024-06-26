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
	
	router.HandlerFunc(http.MethodGet, "/v1/players", app.requirePermission("players:read",app.listPlayersHandler))
	router.HandlerFunc(http.MethodPost, "/v1/players", app.requirePermission("players:write",app.createPlayerHandler))
	router.HandlerFunc(http.MethodGet, "/v1/players/:id",app.requirePermission("players:read", app.showPlayerHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/players/:id",app.requirePermission("players:write",app.updatePlayerHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/players/:id",app.requirePermission("players:write",app.deletePlayerHandler))


	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
