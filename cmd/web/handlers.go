package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.tushar.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	http.NotFound(w, r)
	// 	return
	// }

	snippets, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// for _, snippet := range snippets {
	// 	fmt.Printf("%+v\n", snippet)
	// }

	data := app.newTemplateData(r)
	data.Snippets = snippets

	// helper to render the tmpl-page passed
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// helper method to get the common dynamic data
	data := app.newTemplateData(r)
	// handler or api specific data
	data.Snippet = snippet

	// helper to render the tmpl-page passed.
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		// Changing the response header map after a call to w.WriteHeader() or w.Write() will have no effect
		w.Header().Set("Allow", "POST")
		w.Header()["Date"] = nil // suppressing default system-generated headers in response

		app.clientError(w, http.StatusMethodNotAllowed)
		// w.Header()["Allow"] = []string{"POST"} // direct assignment
		// w.Header().Set("Content-Type", "application/json") // to set the content-type explicity
		// can only be called only once, default value set is 200Ok inside Write(), so we have to set it before calling Write()
		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte("Method Not Allowed"))
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // write the msg and set the response status internally
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

// Caching templates
