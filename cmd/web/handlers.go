package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

// define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
  // Use the Header().Add() method to add a 'Server: Go' header to the
  // response header map. The first parameter is the header name, and
  // the second parameter is the header value.	
	w.Header().Add("Server", "Go")

	// Initialize a slice containing the paths to the two files. It's important
  // to note that the file containing our base template must be the *first*
  // file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use the template.ParseFiles() function to read the files and store the
  // templates in a template set. Notice that we use ... to pass the contents 
  // of the files slice as variadic arguments.
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Then we use the Execute() method on the template set to write the
  // template content as the response body. The last parameter to Execute()
  // represents any dynamic data that we want to pass in, which for now we'll
  // leave as nil.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	// Extract the value of the id wildcard from the request using r.PathValue()
  // and try to convert it to an integer using the strconv.Atoi() function. If
  // it can't be converted to an integer, or the value is less than 1, we
  // return a 404 page not found response.
	var invalid bool = (err != nil || id < 1)
	if invalid {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Sprintf() function to interpolate the id value with a
  // message, then write it as the HTTP response.
	fmt.Fprintf(w, "Display a specific input with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Display a form for creating a new snippet...")
}


func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// use the w.WriteHeader() method to send an 201 status code.
	w.WriteHeader(http.StatusCreated)

	// Then w.Write() method to write the response body as normal.
	io.WriteString(w, "Save a new snippet...")
}
