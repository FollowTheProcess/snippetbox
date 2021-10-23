package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (a *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.WithFields(logrus.Fields{"Addr": r.RemoteAddr, "Method": r.Method, "URL": r.URL.RequestURI()}).Infoln()

		next.ServeHTTP(w, r)
	})
}

func (a *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function which will always run in the event
		// of a panic as Go unwinds the call stack
		defer func() {
			// Check if there's been a panic, and if so, handle it
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				a.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		// Carry on
		next.ServeHTTP(w, r)
	})
}
