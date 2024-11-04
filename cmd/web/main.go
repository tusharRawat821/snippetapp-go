package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	// Go MySQL Driver is an implementation of Go's
	// database/sql/driver interface
	// blank import as we want our driver's init() func
	// to run so that it can register itself with db/sql pkg
	_ "github.com/go-sql-driver/mysql"
	"snippetbox.tushar.net/internal/models"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippetModel  *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// open db connection here
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	// app ever only is terminated by fatal log or ctrl + c and in both cases
	// defer won't run, so this call is a bit flawed, but still a good practice
	// to do so
	defer db.Close()

	// parsing all html-templ files in-memory
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippetModel:  &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	// db obj is a pool of connections, but doesn't actually
	// create any connections, actuall connec are established lazily only.
	// safe for concurrent access and long-lived object. so should be created once only in main ideally
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// to actually see if connec is being made successfully
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
