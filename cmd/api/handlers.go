package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strconv"

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
		books, err := app.models.Books.GetAll()

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := app.writeJSON(w, http.StatusOK, envelope{"books": books}, nil); err != nil {
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

		book := &data.Book{
			Title:     input.Title,
			Published: input.Published,
			Pages:     input.Pages,
			Genres:    input.Genres,
			Rating:    input.Rating,
		}

		err = app.models.Books.Insert(book)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("v1/books/%d", book.ID))

		err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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

	book, err := app.models.Books.Get(idInt)
	if err != nil {
		app.logger.Println(err.Error())
		switch {
		case err.Error() == "record not found":
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updateBook(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)
	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	book, err := app.models.Books.Get(idInt)

	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	var input struct {
		Title     *string  `json:"title"`
		Published *int     `json:"published"`
		Pages     *int     `json:"pages"`
		Genres    []string `json:"genres"`
		Rating    *float64 `json:"rating"`
	}

	err = app.readJSON(w, r, &input)
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

	err = app.models.Books.Update(book)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil); err != nil {
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

	err = app.models.Books.Delete(idInt)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("record not found")):
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
