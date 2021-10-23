package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler {
	// Alice let's you build a composable chain of middleware applied left to right
	standardMiddleware := alice.New(a.recoverPanic, a.logRequest, secureHeaders)

	router := mux.NewRouter()

	// HTTP GET Handlers
	router.HandleFunc("/", a.home).Methods(http.MethodGet)
	router.HandleFunc("/snippet/{id:[0-9]+}", a.showSnippet).Methods(http.MethodGet)

	// HTTP POST Handlers
	router.HandleFunc("/snippet/create", a.createSnippet).Methods(http.MethodPost)

	// Static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/")))).Methods(http.MethodGet)

	// Apply the chain of middleware in the right order and then call our router
	return standardMiddleware.Then(router)
}
