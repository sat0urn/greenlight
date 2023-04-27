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
	// app.requireOnlyAdmin(app.requireActivatedUser())
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	// app.requireActivatedUser()
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)
	// app.requireOnlyAdmin(app.requireActivatedUser())
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.updateMovieHandler)
	// app.requireOnlyAdmin(app.requireActivatedUser())
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.deleteMovieHandler)
	// app.requireActivatedUser()
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)

	// app.requireActivatedUser()
	router.HandlerFunc(http.MethodPost, "/v1/trailers", app.createTrailerHandler)

	//app.requireActivatedUser()
	router.HandlerFunc(http.MethodPost, "/v1/directors", app.createDirectorHandler)
	router.HandlerFunc(http.MethodGet, "/v1/directors", app.requireActivatedUser(app.listDirectorsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// app.authenticate()
	return app.recoverPanic(app.rateLimit(router))
}
