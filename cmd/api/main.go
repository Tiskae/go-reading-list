package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const VERSION = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "db-sdn", os.Getenv("READING_LIST_DB_DSN"), "POSTGRESQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := application{
		config: cfg,
		logger: logger,
	}

	db, errOpenDB := sql.Open("postgres", cfg.dsn)
	if errOpenDB != nil {
		logger.Fatal(errOpenDB)
	}

	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		logger.Fatal(errPing)
	}

	logger.Printf("database connection pool was established")

	addr := fmt.Sprintf(":%d", cfg.port)

	srv := http.Server{
		Addr:         addr,
		Handler:      app.Route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Printf("Server listening on port %d", app.config.port)

	err := srv.ListenAndServe()
	logger.Fatal(err)
}
