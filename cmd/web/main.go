package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	// Import the models package that we just created. You need to prefix this with
	// whatever module path you set up back in chapter 02.01 (Project Setup and Creating
	// a Module) so that the import statement looks like this:
	// "{your-module-path}/internal/models". If you can't remember what module path you
	// used, you can find it at the top of the go.mod file.
	"github.com/mixnblend/snippetbox/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	// Define a new command-line flag witht the name addr, a default value of ":4000"
	// and some short helpt text explaining what the flag controls. The value of the flag
	// will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// define a new command-line flag for the MYSQL DSN string.
	dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL, data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This read in the command-line flag value and assigns it to the addr variable.
	// You need to call this *before* you use the addr variable, otherwise
	// it will always contain the default value of ":4000". If any errors are encountered
	// during parsing the application will be terminated.
	flag.Parse()

	// Use he slog.New() function to initialise a new structured logger, which
	// writes to the standard out and uses the default settings.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() function below. We pass openDB() the DSN
	// from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	// Initialise a new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a new instance of our application struct, containing the
	// dependencies (for now, just the structured logger).
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it (i.e. prefix it with
	// the * symbol) before using it. Note that we're using the log.Printf()
	// function to interpolate the address with the log message.
	logger.Info("starting server", slog.String("addr", *addr))

	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	err = http.ListenAndServe(*addr, app.routes())

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
