package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FollowTheProcess/snippetbox/pkg/forms"
	"github.com/FollowTheProcess/snippetbox/pkg/models"
	"github.com/gorilla/mux"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.render(w, r, "home.page.tmpl", &templateData{Snippets: s})
}

func (a *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// parse the snippet id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		a.serverError(w, err)
	}

	s, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}
		return
	}

	a.render(w, r, "show.page.tmpl", &templateData{Snippet: s})
}

func (a *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		a.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := a.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (a *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, "create.page.tmpl", nil)
}
