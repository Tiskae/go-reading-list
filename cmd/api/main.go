package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const VERSION = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8081, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthCheck)

	addr := fmt.Sprintf(":%d", cfg.port)

	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println(err)
	}
}
