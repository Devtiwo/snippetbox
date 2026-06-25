package main

import (
  "net/http"
  "github.com/justinas/alice"
  "github.com/julienschmidt/httprouter"
)

// The routes() method returns a servemux containing the application routes.
// Updating the signature for the routes() method so that it returns a http.Handler instead of *http.ServeMux
func (app *application) routes() http.Handler {
  router := httprouter.New()
  router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    app.notFound(w)
  })
  fileServer := http.FileServer(http.Dir("./ui/static/"))

  // Update the pattern for the route for the static files.
  // Using the http.StripPrefix() function to remove the "/static" prefix from the request URL path before passing it to the file server.
  router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
  router.HandlerFunc(http.MethodGet, "/", app.home)
  router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
  router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
  router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

  // Creating a middleware chain containing our standard middleware which will be used for every request our application receives.
  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  // Wrap the router with the middleware and return it as normal.
  return standard.Then(router)
}