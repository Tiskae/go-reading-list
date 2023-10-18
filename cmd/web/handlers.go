package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "The home page")
}

func (app *application) BookView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "View a single page")
}

func (app *application) bookCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new book record form")
}
