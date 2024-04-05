package main

import (
	"fmt"
	"io"
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

	io.WriteString(w, "Hello from Snippetbox")
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
