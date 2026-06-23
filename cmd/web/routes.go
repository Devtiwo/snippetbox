package main

import (
  "net/http"
  "github.com/justinas/alice"
)

// The routes() method returns a servemux containing the application routes.
// Updating the signature for the routes() method so that it returns a http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
  mux := http.NewServeMux()
  fileServer := http.FileServer(http.Dir("./ui/static/"))

  // Using the mux.Handle() method to register the file server handler for the "/static/" route.
  // Using the http.StripPrefix() function to remove the "/static" prefix from the request URL path before passing it to the file server.
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  // Registering the handler functions for different routes.
  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/snippet/view", app.snippetView)
  mux.HandleFunc("/snippet/create", app.snippetCreate)

  // Creating a middleware chain containing our standard middleware which will be used for every request our application receives.
  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  // Returning the 'standard' middleware chain followed by the servemux.
  return standard.Then(mux)
}