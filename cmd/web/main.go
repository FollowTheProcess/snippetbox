package main

import (
	"context"
	"database/sql"
	"flag"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FollowTheProcess/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = ":8000"
	defaultDSN  = ""
)

type application struct {
	logger        *logrus.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Accept command line flags for configuration and secrets
	port := flag.String("port", defaultPort, "HTTP network address")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	flag.Parse()

	// Set up logger
	log := logrus.New()
	log.Out = os.Stdout

	if *dsn == "" {
		log.Fatalln("dsn must not be empty")
	}

	log.Infoln("Establishing db connection")
	db, err := openDB(*dsn)
	if err != nil {
		log.WithError(err).Fatalln("Error connecting to DB")
	}
	defer func() {
		log.Infoln("Closing DB connection")
		db.Close()
	}()

	// Initialise a new template cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		log.WithError(err).Fatalln("Error initialising template cache")
	}

	app := &application{
		logger:        log,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,  // max time to read request from the client
		WriteTimeout: 10 * time.Second, // max time to write response to the client
		IdleTimeout:  60 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server in a goroutine so it runs off doing it's own thing
	go func() {
		log.WithField("port", *port).Infoln("Starting server on port")

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.WithError(err).Errorln("Error starting server")
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block the rest of the code until a signal is received.
	sig := <-c
	log.WithField("sig", sig).Infoln("Got signal")
	log.Infoln("Shutting everything down gracefully")

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Fatalln("Graceful shutdown failed")
	}
	log.Infoln("Server shutdown successfully")
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
