package main

import (
  "fmt"
  "errors"
  "net/http"
  "strconv"
  "github.com/Devtiwo/snippetbox/internal/models"
)

// Creates a home page handler function.
// Changing the structure of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
	app.notFound(w) // Using the notFound helper method to send a 404 Not Found response to the user if the URL path is not "/".
	return
  }
  snippets, err := app.snippets.Latest()
  if err != nil {
    app.serverError(w, err)
   return
  }
  // Calling the newTemplateData() helper to get a templateData struct containing the default data
  data := app.newTemplateData(r)
  data.Snippets = snippets
  // Using the new render helper
  app.render(w, http.StatusOK, "home.tmpl", data)
}

// Creates a snippetview handler function.
// Changing the structure of the snippetView handler so it is defined as a method against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
  id, err:= strconv.Atoi(r.URL.Query().Get("id"))
  if err != nil || id < 1 {
	app.notFound(w) // Using the notFound helper.
	return
  }
  snippet, err := app.snippets.Get(id)
  if err != nil {
	if errors.Is(err, models.ErrNoRecord) {
	  app.notFound(w)
	} else {
	  app.serverError(w, err)
	}
	return
  }
  data := app.newTemplateData(r)
  data.Snippet = snippet
  // Using the new render helper
  app.render(w, http.StatusOK, "view.tmpl", data)
}

// Creates a snippetCreate handler function.
// Changing the structure of the snippetCreate handler so it is defined as a method against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
	w.Header().Set("Allow", http.MethodPost)
	app.clientError(w, http.StatusMethodNotAllowed) // Using the clientError helper
	return
  }
  title := "O snail"
  content := "O snail\nClimb Mount Fuji, \nBut slowly, slowly!\n\n- Kobayashi Issa"
  expires := 7
  id, err := app.snippets.Insert(title, content, expires)
  if err != nil {
	app.serverError(w, err)
	return
  }
  http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
  w.Write([]byte("Creates a new snippet..."))
}