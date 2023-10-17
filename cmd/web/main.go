package main

import (
	"flag"
	"log"
	"net/http"
)

type application struct {
}

func main() {
	app := application{}

	addr := flag.String("address", ":80", "HTTP network address")

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
	}

	log.Printf("Listening on port %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
