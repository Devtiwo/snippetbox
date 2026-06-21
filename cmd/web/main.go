package main

import (
  "os"
  "flag"
  "log"
  "net/http"
)

func main() {
    // Using the flag package to define a command-line flag named "addr" with a default value of ":4000" and a description of "HTTP network address".
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Calling the flag.Parse() function to parse the command-line flags and store the values in the corresponding variables.
	flag.Parse()

	// Using log.New() function to create a new logger for writing information messages. It takes three parameters.
	// Destination to write the logs to (os.Stdout), a string prefix for the message, and flags to indicate what additional information to include which are joined using the bitwise OR operator.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but usestderr as the destination and use the log.Lshortfile flag to include the relevant file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    mux := http.NewServeMux()
    
	// Creating a file server to serve static files from the "./ui/static/" directory.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

    // Using the mux.Handle() method to register the file server handler for the "/static/" route.
    // Using the http.StripPrefix() function to remove the "/static" prefix from the request URL path before passing it to the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Registering the handler functions for different routes.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	
	infoLog.Printf("starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
	  errorLog.Fatal(err)
	}
}