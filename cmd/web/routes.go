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

	// 	For the sessions to work, we also need to wrap our application routes with the middleware
	// provided by the SessionManager.LoadAndSave() method. This middleware automatically
	// loads and saves session data with every HTTP request and response
	// Create a new middleware chain containing the middleware specific to our
	// dynamic application routes. For now, this chain will only contain the
	// LoadAndSave session middleware but we'll add more to it later.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// mux := http.NewServeMux()					                              // This is a middleware handler which keeps a map of {path : handler} and does the re-direction
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))                             // exact match to "/{$}" path
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))      // fixed path and not a subtree
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))      // get create snippet form
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost)) // save snippet

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
