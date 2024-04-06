package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
    
    // Define a new command-line flag witht the name addr, a default value of ":4000"
    // and some short helpt text explaining what the flag controls. The value of the flag
    // will be stored in the addr variable at runtime.
    addr := flag.String("addr", ":4000", "HTTP network address")
    
    // Importantly, we use the flag.Parse() function to parse the command-line flag.
    // This read in the command-line flag value and assigns it to the addr variable.
    // You need to call this *before* you use the addr variable, otherwise
    // it will always contain the default value of ":4000". If any errors are encountered
    // during parsing the application will be terminated.
    flag.Parse()

    // Use he slog.New() function to initialise a new structured logger, which
    // writes to the standard out and uses the default settings.
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", home)
    mux.HandleFunc("GET /snippet/view/{id}", snippetView)
    mux.HandleFunc("GET /snippet/create", snippetCreate)
    mux.HandleFunc("POST /snippet/create", snippetCreatePost)

    // The value returned from the flag.String() function is a pointer to the flag
    // value, not the value itself. So in this code, that means the addr variable 
    // is actually a pointer, and we need to dereference it (i.e. prefix it with
    // the * symbol) before using it. Note that we're using the log.Printf() 
    // function to interpolate the address with the log message.
    logger.Info("starting server", slog.String("addr", *addr))
    
    // And we pass the dereferenced addr pointer to http.ListenAndServe() too.
    err := http.ListenAndServe(*addr, mux)
    // And we also use the Error() method to log any error message returned by
    // http.ListenAndServe() at Error severity (with no additional attributes),
    // and then call os.Exit(1) to terminate the application with exit code 1.
    logger.Error(err.Error())
    os.Exit(1)
}
