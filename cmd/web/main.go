package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore" // New import
	"github.com/alexedwards/scs/v2"         // New import
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mixnblend/snippetbox/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include the structured logger, but we'll
// add more to this as the build progresses.
type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	// Initialize a decoder instance ...
	formDecoder := form.NewDecoder()

	// Use the scs.New() function to initialize a new session manager. Then we
	// configure it to use our MySQL database as the session store, and set a
	// lifetime of 12 hours (so that sessions automatically expire 12 hours
	// after first being created).
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// And add the session manager to our application dependencies.
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initalise a new http.Server struct. We set the Addr and Handler fields
	// so that the server uses the same network address and routes as before.

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Create a *log.Logger from our structured logger handler, which writes
		// log entries at Error level, and assign it to the ErrorLog field. If
		// you would prefer to log the server errors at Warn level instead, you
		// could pass slog.LevelWarn as the final parameter.
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So in this code, that means the addr variable
	// is actually a pointer, and we need to dereference it (i.e. prefix it with
	// the * symbol) before using it. Note that we're using the log.Printf()
	// function to interpolate the address with the log message.
	logger.Info("starting server", slog.String("addr", *addr))

	// And we pass the dereferenced addr pointer to http.ListenAndServe() too.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

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
