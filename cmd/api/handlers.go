package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"time"

	"github.com/tiskae/go-reading-list/internal/data"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	data := map[string]string{
		"status":     "available",
		"envionment": app.config.env,
		"version":    VERSION,
	}

	js, err := json.Marshal(data)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}

func (app *application) getCreateBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// GET -> /v1/books
		books := []data.Book{
			{
				ID:        8761890,
				CreatedAt: time.Now(),
				Title:     "So Long a Letter",
				Published: 2005,
				Pages:     230,
				Genres:    []string{"Epistolary", "fiction"},
				Rating:    3.9,
				Version:   1,
			},
			{
				ID:        1678921,
				CreatedAt: time.Now(),
				Title:     "Greengold Autumn",
				Published: 1971,
				Pages:     171,
				Genres:    []string{"Romance"},
				Rating:    4.5,
				Version:   1,
			},
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"books": books}); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	} else if r.Method == http.MethodPost {
		// POST -> /v1/books
		var input struct {
			Title     string   `json:"title"`
			Published int      `json:"published"`
			Pages     int      `json:"pages"`
			Genres    []string `json:"genres"`
			Rating    float64  `json:"rating"`
		}

		err := app.readJSON(w, r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "%v\n", input)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeleteBooksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// GET -> /v1/books/{id}
		app.getBook(w, r)

	case http.MethodPut:
		// PUT -> /v1/books/{id}
		app.updateBook(w, r)

	case http.MethodDelete:
		// DELETE -> /v1/books/{id}
		app.deleteBook(w, r)

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getBook(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	book := data.Book{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Gifted Hands",
		Published: 1990,
		Pages:     180,
		Genres:    []string{"Motivational", "Self development", "Autobiography"},
		Rating:    4.7,
		Version:   2,
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	idInt, errParse := strconv.ParseInt(id, 10, 64)

	if errParse != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	var input struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float64 `json:"rating"`
	}

	book := data.Book{
		ID:        idInt,
		CreatedAt: time.Now(),
		Title:     "Letter to Zoyah ❤️",
		Published: 2023,
		Pages:     2,
		Genres:    []string{"Love"},
		Rating:    5,
		Version:   1,
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if input.Title != nil {
		book.Title = *input.Title
	}

	if input.Published != nil {
		book.Published = *input.Published
	}

	if input.Pages != nil {
		book.Pages = *input.Pages
	}

	if input.Rating != nil {
		book.Rating = *input.Rating
	}

	if len(input.Genres) > 0 {
		book.Genres = input.Genres
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteBook(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	fmt.Fprintf(w, "Delete book with id %d", idInt)
}
