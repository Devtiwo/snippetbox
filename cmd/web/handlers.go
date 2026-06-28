package main

import (
  "fmt"
  "errors"
  "net/http"
  "strconv"
  "github.com/Devtiwo/snippetbox/internal/models"
  "github.com/Devtiwo/snippetbox/internal/validator"
  "github.com/julienschmidt/httprouter"
)

// Define a snippetCreateForm struct to represent the form data and validation errors for the form fields
// Update our snippetCreateForm struct to include struct tags which tell the decoder how to map HTML form values into the different struct fields.
// So, for example, here we're telling the decoder to store the value from the HTML form input with the name "title" in the Title field.
// The struct tag `form:"-"` tells the decoder to completely ignore a field during decoding
type snippetCreateForm struct {
  Title string `form:"title"`
  Content string `form:"content"`
  Expires int     `form:"expires"`
  validator.Validator `form:"-"`
}
// Creates a home page handler function.
// Changing the structure of the home handler so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
  // When httprouter is parsing a request, the values of any named parameters will be stored in the request contest which we can then use ParamsFromContext() function to retrieve a slice containing these parameter names and values.
  params := httprouter.ParamsFromContext(r.Context())

  // We can then use the ByName() method to get the value of the "id" named parameter from the slice and validate it as normal.
  id, err:= strconv.Atoi(params.ByName("id"))
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

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
  data := app.newTemplateData(r)
  // Initialize a new createSnippetForm instance and pass it to the template. This is also an opportunity to set default values for the form.
  data.Form = snippetCreateForm{
    Expires: 365,
  }
  app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
  // Declare a new empty instance of the snippetCreateForm struct.
  var form snippetCreateForm

  // Call the Decode() method of the form decoder, passing in the current request and *a pointer* to our snippetCreateForm struct.
  // This will essentially fill our struct with the relevant values from the HTML form.
  // If there is a problem, we return a 400 Bad Request response to the client
  err := app.decodePostForm(r, &form)
  if err != nil {
    app.clientError(w, http.StatusBadRequest)
    return
  }

  // calling CheckField() directly on snippetCreateForm struct to execute our validation checks
  form.CheckField(validator.NotBlank(form.Title), "title", "this field cannot be blank")
  form.CheckField(validator.MaxChars(form.Title, 100), "title", "this field cannot be more than 100 characters long")
  form.CheckField(validator.NotBlank(form.Content), "content", "this field cannot be blank")
  form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "this field must equal 1, 7 or 365")

  // Use the Valid() method to see if any of the checks failed. If they did, then re-render the template passing in the form in the same way as before
  if !form.Valid() {
    data := app.newTemplateData(r)
    data.Form = form
    app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
    return
  }
  id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
  if err != nil {
    app.serverError(w, err)
    return
  }
  
  // Using the Put() method to add a string value ("Snippet successfully created!") and the corresponding key ("flash") to the session data.
  app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
  http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}