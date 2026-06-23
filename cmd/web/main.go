package main

import (
  "os"
  "flag"
  "log"
  "net/http"
  "html/template"
  "database/sql"
  "github.com/Devtiwo/snippetbox/internal/models"
  _ "github.com/go-sql-driver/mysql"
)

// Defining an application struct to hold the dependencies for our web app.
type application struct {
  errorLog *log.Logger
  infoLog *log.Logger
  snippets *models.SnippetModel
  templateCache map[string]*template.Template 
}

func main() {
    // Using the flag package to define a command-line flag named "addr" with a default value of ":4000" and a description of "HTTP network address".
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Using the flag package to define a command-line flag for mysql dsn string.
	dsn := flag.String("dsn", "web:adewunmi5020@/snippetbox?parseTime=true", "MySQL data source name")

	// Calling the flag.Parse() function to parse the command-line flags and store the values in the corresponding variables.
	flag.Parse()

	// Using log.New() function to create a new logger for writing information messages. It takes three parameters.
	// Destination to write the logs to (os.Stdout), a string prefix for the message, and flags to indicate what additional information to include which are joined using the bitwise OR operator.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but usestderr as the destination and use the log.Lshortfile flag to include the relevant file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    
	// Calling the openDB() function to create a connection pool for the MySQL database.
	db, err := openDB(*dsn)
	if err != nil {
	  errorLog.Fatal(err)
	}

    // Using the defer statement to ensure that the database connection pool is closed before the main() function exits.
	defer db.Close()

	// initializing a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
	  errorLog.Fatal(err)
	}
	// initializing a new instance of the application struct containing the dependencies.
	app := &application{
	  errorLog: errorLog,
	  infoLog: infoLog,
	  snippets: &models.SnippetModel{DB: db},
	  templateCache: templateCache,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so that the server uses the same network address and routes before and set the errorlog field so that the server now uses the custome errorlog logger in the event of any problems.
	srv := &http.Server{
	  Addr: *addr,
	  ErrorLog: errorLog,
	  Handler: app.routes(), // Calling the new app.routes() method to get the servemux containing our routes.
	}

	infoLog.Printf("starting server on %s", *addr)

	// Call the ListenAndServe() method on our new http.Server struct.
	// Since the err variable is already declared, we can use the assignment operator to assign the return value of the ListenAndServe() method to it.
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool for a given DSN.
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