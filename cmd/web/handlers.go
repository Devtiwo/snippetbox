package main

import (
  "fmt"
  "log"
  "net/http"
  "strconv"
  "html/template"
)

// Creates a home page handler function.
func home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
	http.NotFound(w, r)
	return
  }

  // Initailizing a slice containing the paths of the template files to be parsed
  // The base file must come first, followed by the other template files.
  files := []string {
	"./ui/html/base.tmpl",
	"./ui/html/partials/nav.tmpl",
	"./ui/html/pages/home.tmpl",
  }

  // using the template.ParseFiles() function to read the template file and storing it into a template set.
  ts, err := template.ParseFiles(files...)
  if err != nil {
	log.Print(err.Error())
	http.Error(w, "Internal server error", 500)
    return
  }
  
  // Using the ExecuteTemplate() method on the template set to write the template content as the response body.
  // The second parameter in the Execute() methos is used to pass dynamic data to the template.
  err = ts.ExecuteTemplate(w, "base", nil)
  if err != nil {
	log.Print(err.Error())
	http.Error(w, "Internal server error", 500)
  }}

// Creates a snippetview handler function.
func snippetView(w http.ResponseWriter, r *http.Request) {
  id, err:= strconv.Atoi(r.URL.Query().Get("id"))
  if err != nil || id < 1 {
	http.NotFound(w, r)
	return
  }
  fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Creates a snippetCreate handler function.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
	w.Header().Set("Allow", http.MethodPost)
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
  }
  w.Write([]byte("Creates a new snippet..."))
}