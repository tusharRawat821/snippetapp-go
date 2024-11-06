package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// file server = serving in http response
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// mux := http.NewServeMux()                                                // This is a middleware handler which keeps a map of {path : handler} and does the re-direction
	router.HandlerFunc(http.MethodGet, "/", app.home)                             // exact match to "/{$}" path
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)      // fixed path and not a subtree
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)      // get create snippet form
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost) // save snippet

	// composable middleware and cleanr/easier to understand using alice pkg
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}

// 3 pl routers -  julienschmidt/httprouter, go-chi/chi and gorilla/mux

// using httprouter
// no conflicting routes exist
// named param route - GET /foo/:name
// catch all param after the path route -  GET /foo/*name
// eg : router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
