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
	"github.com/tiskae/go-reading-list/internal/data"
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
	models data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.dsn, "db-sdn", os.Getenv("READINGLIST_DB_DSN"), "POSTGRESQL DSN")
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
	// postgresql://postgres:62135665@localhost/readinglist?sslmode=disable

	defer db.Close()

	errPing := db.Ping()
	if errPing != nil {
		logger.Fatal(errPing)
	}

	app.models = data.NewModels(db)

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
