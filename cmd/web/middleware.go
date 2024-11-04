package main

import (
	"fmt"
	"net/http"
)

// middleware - headers, logging, authentication, etc.

func (app *application) logRequest(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// set respone headers for all requests
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})

}

func (app *application) recoverPanic(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				// Setting the Connection: Close header on the response acts as a trigger to make Go’s
				// HTTP server automatically close the current connection after a response has been sent. It
				// also informs the user that the connection will be closed
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Extra - Panic recovery in other background goroutines
// So, if you are spinning up additional goroutines from within your web application and there is
// any chance of a panic, you must make sure that you recover any panics from within those too to stop completed app being crashed
// eg:-
// func myHandler(w http.ResponseWriter, r *http.Request) {
// 	...
// 	// Spin up a new goroutine to do some background processing.
// 	go func() {
// 	defer func() {
// 	if err := recover(); err != nil {
// 	log.Print(fmt.Errorf("%s\n%s", err, debug.Stack()))
// 	}
// 	}()
// 	doSomeBackgroundProcessing()
// 	}()
// 	w.Write([]byte("OK"))
// 	}

// composable middleware
