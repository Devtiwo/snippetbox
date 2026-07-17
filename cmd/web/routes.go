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
  
  // Create a new middleware chain containing the middleware specific to our dynamic application routes.
  dynamic := alice.New(app.sessionManager.LoadAndSave)

  // Update these routes to use the new dynamic middleware chain followed by the appropriate handler function. Note that because the alice ThenFunc() method returns a http.Handler (rather than a http.HandlerFunc)
  // We also need to switch to registering the route using the router.Handler() method.
  // Update the pattern for the route for the static files.
  // Using the http.StripPrefix() function to remove the "/static" prefix from the request URL path before passing it to the file server.
  router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
  router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
  router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
  router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
  router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
  router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
  router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
  router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
  router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
  router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

  // Creating a middleware chain containing our standard middleware which will be used for every request our application receives.
  standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
  // Wrap the router with the middleware and return it as normal.
  return standard.Then(router)
}