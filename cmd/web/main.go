package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/FollowTheProcess/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

const (
	defaultPort = ":8000"
	defaultDSN  = "web:snippetpassword@/snippetbox?parseTime=true"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	addr := flag.String("addr", defaultPort, "HTTP network address")
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Println("Establishing db connection")
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer func() {
		infoLog.Println("Closing DB connection")
		db.Close()
	}()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server in a goroutine so it runs off doing it's own thing
	go func() {
		infoLog.Printf("Starting server on %s\n", *addr)

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errorLog.Printf("Error starting server: %s\n", err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block the rest of the code until a signal is received.
	sig := <-c
	infoLog.Println("Got signal:", sig)
	infoLog.Println("Shutting everything down gracefully")

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		errorLog.Fatal("Graceful shutdown failed")
	}
	infoLog.Println("Server shutdown successfully")
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
