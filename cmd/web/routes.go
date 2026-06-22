package main

import "net/http"

// The routes() method returns a servemux containing the application routes.
func (app *application) routes() *http.ServeMux {
  mux := http.NewServeMux()
  fileServer := http.FileServer(http.Dir("./ui/static/"))

  // Using the mux.Handle() method to register the file server handler for the "/static/" route.
  // Using the http.StripPrefix() function to remove the "/static" prefix from the request URL path before passing it to the file server.
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  // Registering the handler functions for different routes.
  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/snippet/view", app.snippetView)
  mux.HandleFunc("/snippet/create", app.snippetCreate)
  return mux
}