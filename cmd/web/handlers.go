package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.tushar.net/internal/constants"
	"snippetbox.tushar.net/internal/validator"
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

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 0 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, constants.ErrNoRecord) {
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

	data := app.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
	// w.Write([]byte("Display the form for reating a snippet..."))
}

// represent the form data entered + validator
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expires int    `form:"expires"`
	// struct embedding : re-usability with composition
	// embedding the struct inside another struct
	validator.Validator `form:"-"` // struct tag `form:"-"` used to tell decoder to ignore field during decoding
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// if r.Method != http.MethodPost {
	// Changing the response header map after a call to w.WriteHeader() or w.Write() will have no effect
	// w.Header().Set("Allow", "POST")
	// w.Header()["Date"] = nil // suppressing default system-generated headers in response

	// app.clientError(w, http.StatusMethodNotAllowed)
	// w.Header()["Allow"] = []string{"POST"} // direct assignment
	// w.Header().Set("Content-Type", "application/json") // to set the content-type explicity
	// can only be called only once, default value set is 200Ok inside Write(), so we have to set it before calling Write()
	// w.WriteHeader(http.StatusMethodNotAllowed)
	// w.Write([]byte("Method Not Allowed"))
	// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // write the msg and set the response status internally
	// return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippetModel.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// snippet is created successfully in db
	// then we can store data in the session with key = flash
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
