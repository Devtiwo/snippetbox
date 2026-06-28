package main

import (
  "html/template"
  "path/filepath"
  "time"
  "github.com/Devtiwo/snippetbox/internal/models"
)

type templateData struct {
  CurrentYear int
  Snippet *models.Snippet
  Snippets []*models.Snippet
  Form any
  Flash string
}

// Creates a humanDate which returns a nicely formatted string representation of the time.Time object.
func humanDate(t time.Time) string{
  return t.Format("02 Jan 2026 at 03:21 am")
}

// Initialize a template.FuncMap object and store it in a global variable.
var functions = template.FuncMap{
  "humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
  // initializing a new map to act as cache
  cache := map[string]*template.Template{}

  // using the filepath.Glob() function to get a slice of all file path that match the pattern "./ui/html/pages/*.tmpl".
  // This gives us a slice of all the filepath for our app page template eg [ui/html/pages/home.tmpl ui/html/pages/view.tmpl].
  pages, err := filepath.Glob("./ui/html/pages/*tmpl")
    if err != nil {
	  return nil, err
	}

  // looping through the page filepath one after the other
  for _, page := range pages {
	// Extract the file name from the full filepath eg (home.tmpl) and assign it to the name variable
	name := filepath.Base(page)

	// the template.FuncMap must be registered with the template set before the ParseFiles() method is called.
  // This means we have to use template.New() to create an empty template set, use the Func() method to register the template.FuncMap, and then parse the file as normal. 
	ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
	if err != nil {
	  return nil, err
	}
    // Call ParseGlob() on this template set to add any partials
	ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
	if err != nil {
	  return nil, err
	}
	// call ParseFiles on this template set to add the page
	ts, err = ts.ParseFiles(page)
	if err != nil {
	  return nil, err
	}
	// Add the template set to the map using the name of page eg (home.tmpl) as the key
	cache[name] = ts
  }
  // Return the map
  return cache, nil
  }