package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()                        // This is a middleware handler which keeps a map of {path : handler} and does the re-direction
	mux.HandleFunc("/{$}", app.home)                 // exact match to "/" path
	mux.HandleFunc("/snippet/view", app.snippetView) // fixed path and not a subtree
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// file server = serving in http response
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// composable middleware and cleanr/easier to understand
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}
