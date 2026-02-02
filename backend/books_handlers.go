package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type BooksAPI struct {
	store *BookStore
}

func NewBooksAPI(store *BookStore) *BooksAPI {
	return &BooksAPI{store: store}
}

// GetBooksHandler godoc
// @Summary List all books
// @Tags books
// @Produce json
// @Success 200 {array} Book
// @Failure 500 {object} errorResponse
// @Router /books [get]
func (api *BooksAPI) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	books, err := api.store.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, books)
}

// CreateBookHandler godoc
// @Summary Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body Book true "Book"
// @Success 201 {object} Book
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /books [post]
func (api *BooksAPI) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	created, err := api.store.Create(r.Context(), b)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// GetBookHandler godoc
// @Summary Get a book by ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} Book
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /books/{id} [get]
func (api *BooksAPI) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	b, err := api.store.Get(r.Context(), id)
	if err == ErrNotFound {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "book not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, b)
}

// UpdateBookHandler godoc
// @Summary Update a book by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body Book true "Book"
// @Success 200 {object} Book
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /books/{id} [put]
func (api *BooksAPI) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}

	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	updated, err := api.store.Update(r.Context(), id, b)
	if err == ErrNotFound {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "book not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DeleteBookHandler godoc
// @Summary Delete a book by ID
// @Tags books
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /books/{id} [delete]
func (api *BooksAPI) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(w, r, "id")
	if !ok {
		return
	}
	err := api.store.Delete(r.Context(), id)
	if err == ErrNotFound {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "book not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseIDParam(w http.ResponseWriter, r *http.Request, name string) (int64, bool) {
	raw := chi.URLParam(r, name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid id"})
		return 0, false
	}
	return id, true
}
