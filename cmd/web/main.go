package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luckyparakh/snippetbox/pkg/models/mysql"
)

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "address to listen")
	dsn := flag.String("dsn", "web:web@/snippetbox?parseTime=true", "MYSQL data source name")
	flag.Parse()
	db, err := openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		log.Fatal(err)
	}
	app := &application{
		infoLog:       log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile),
		errLog:        log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	app.infoLog.Printf("Starting server on %v\n", *addr)
	// err:= http.ListenAndServe(*addr, mux)
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: app.errLog,
		Handler:  app.routes(),
	}
	err = srv.ListenAndServe()
	app.errLog.Fatal(err)
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
