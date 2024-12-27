package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/stdlib"

	"github.com/mf751f/snippetbox/internal/models"
)

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String(
		"dsn",
		"postgres://postgres:1319@localhost:5432/snippetbox",
		"Postgresql data source name",
	)
	flag.Parse()

	infoLog := log.New(os.Stdout, "[info ]\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stdout, "[Error]\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err)
	}

	app := application{
		infoLog:       infoLog,
		errLog:        errLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %v", *addr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	sql.Register("postgres", stdlib.GetDefaultDriver())
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
