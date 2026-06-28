package main

import (
  "fmt"
  "bytes"
  "errors"
  "net/http"
  "time"
  "runtime/debug"
  "github.com/go-playground/form/v4"
)

// The serverError helper writes an error message and stack trace to the errorLog, then sends a generic 500 internal server error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
  trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
  app.errorLog.Output(2,trace)
  http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
} 

// The clientError helper sends a specific status code and corresponding description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
  http.Error(w, http.StatusText(status), status)
}

// For consistency, we can also implement a notFound helper which sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
  app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
  // Retrieve the appropriate template set from the cache based on the page name ge(home.tmpl) and if no entry exist in the cache with the provided name, create a new error and call the serverError() helper.
  ts, ok := app.templateCache[page]
  if !ok {
    err := fmt.Errorf("the template %s does not exixt", page)
    app.serverError(w, err)
  }

  // initializing a new buffer
  buf := new(bytes.Buffer)
  // writing the template to the buffer first instead of the http.ResponseWriter and if there is any error, it should call the serverError helper and return
  err := ts.ExecuteTemplate(buf, "base", data)
  if err != nil {
    app.serverError(w, err)
    return
  }
  // if there are no errors, then we Write out the http status code to the http.ResponseWriter
  w.WriteHeader(status)
  //write the content of the buffer to the http.ResponseWriter
  buf.WriteTo(w)
}

// Creating a newTemplateData helper that returns a pointer to a templateData struct initialized with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData{
  return &templateData{
    CurrentYear: time.Now().Year(),
    // Add the flash message to the template data, if one exists.
    Flash: app.sessionManager.PopString(r.Context(), "flash"),
  }
}

// Create a new decodePostForm() helper method. The second parameter here, dst, is the target destination that we want to decode the form data into.
func(app *application) decodePostForm(r *http.Request, dst any) error {
  // Call ParseForm() on the request, in the same way that we did in our createSnippetPost handler.
  err := r.ParseForm()
  if err != nil {
    return err
  }
  // Call Decode() on our decoder instance, passing the target destination as the first parameter.
  err = app.formDecoder.Decode(dst, r.PostForm)
  if err != nil {
    // If we try to use an invalid target destination, the Decode() method will return an error with the type *form.InvalidDecoderError.
    // We use errors.As() to check for this and raise a panic rather than returning the error.
    var invalidDecoderError *form.InvalidDecoderError
    if errors.As(err, &invalidDecoderError) {
      panic(err)
    }
    // For all other errors, we return them as normal.
    return err
  }
  return nil
}