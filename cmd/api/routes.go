package main

import "net/http"

func (app *application) Route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthCheck)
	mux.HandleFunc("/v1/books", app.getCreateBooksHandler)
	mux.HandleFunc("/v1/books/", app.getUpdateDeleteBooksHandler)

	return mux
}
