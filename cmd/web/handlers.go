package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	title := "O snail"
	content := "O snail\nClimt Mount Fuji, \nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := "7"

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
