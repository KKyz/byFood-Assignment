package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type errResp struct {
	Error string `json:"error"`
}

func setupTestRouter(t *testing.T) (*chi.Mux, *sql.DB) {
	t.Helper()

	db, err := OpenDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	if err := Migrate(db); err != nil {
		t.Fatal(err)
	}

	store := NewBookStore(db)
	api := NewBooksAPI(store)

	r := chi.NewRouter()

	// Books routes (impt!)
	r.Route("/books", func(r chi.Router) {
		r.Get("/", api.GetBooksHandler)
		r.Post("/", api.CreateBookHandler)
		r.Get("/{id}", api.GetBookHandler)
		r.Put("/{id}", api.UpdateBookHandler)
		r.Delete("/{id}", api.DeleteBookHandler)
	})

	r.Post("/process-url", ProcessURLHandler)

	return r, db
}

func doJSON(t *testing.T, r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func decodeJSON[T any](t *testing.T, rr *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rr.Body.Bytes(), &v); err != nil {
		t.Fatalf("invalid json: %v body=%s", err, rr.Body.String())
	}
	return v
}

func TestBooks_HappyPathCRUD(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// Create
	rr := doJSON(t, r, http.MethodPost, "/books", `{"title":"Dune","author":"Frank Herbert","year":1965}`)
	if rr.Code != http.StatusCreated {
		t.Fatalf("create status got %d body=%s", rr.Code, rr.Body.String())
	}
	created := decodeJSON[Book](t, rr)
	if created.ID <= 0 {
		t.Fatalf("expected ID > 0 got %d", created.ID)
	}

	rr = doJSON(t, r, http.MethodGet, "/books", ``)
	if rr.Code != http.StatusOK {
		t.Fatalf("list status got %d body=%s", rr.Code, rr.Body.String())
	}
	list := decodeJSON[[]Book](t, rr)
	if len(list) != 1 {
		t.Fatalf("expected 1 book, got %d", len(list))
	}

	getPath := fmt.Sprintf("/books/%d", created.ID)
	rr = doJSON(t, r, http.MethodGet, getPath, ``)
	if rr.Code != http.StatusOK {
		t.Fatalf("get status got %d body=%s", rr.Code, rr.Body.String())
	}
	got := decodeJSON[Book](t, rr)
	if got.Title != "Dune" || got.Author != "Frank Herbert" || got.Year != 1965 {
		t.Fatalf("unexpected book: %+v", got)
	}

	rr = doJSON(t, r, http.MethodPut, getPath, `{"title":"Dune Messiah","author":"Frank Herbert","year":1969}`)
	if rr.Code != http.StatusOK {
		t.Fatalf("update status got %d body=%s", rr.Code, rr.Body.String())
	}
	updated := decodeJSON[Book](t, rr)
	if updated.Title != "Dune Messiah" || updated.Year != 1969 {
		t.Fatalf("update mismatch: %+v", updated)
	}

	rr = doJSON(t, r, http.MethodGet, getPath, ``)
	if rr.Code != http.StatusOK {
		t.Fatalf("get-after-update status got %d body=%s", rr.Code, rr.Body.String())
	}
	got2 := decodeJSON[Book](t, rr)
	if got2.Title != "Dune Messiah" || got2.Year != 1969 {
		t.Fatalf("get-after-update mismatch: %+v", got2)
	}

	req := httptest.NewRequest(http.MethodDelete, getPath, nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusNoContent {
		t.Fatalf("delete status got %d body=%s", rr2.Code, rr2.Body.String())
	}

	rr = doJSON(t, r, http.MethodGet, getPath, ``)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("get-after-delete status got %d body=%s", rr.Code, rr.Body.String())
	}
}

func TestBooks_CreateValidationFailures(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	cases := []struct {
		name      string
		body      string
		wantError string
	}{
		{"missing title", `{"author":"A","year":2000}`, "title is required"},
		{"missing author", `{"title":"T","year":2000}`, "author is required"},
		{"year zero", `{"title":"T","author":"A","year":0}`, "year must be > 0"},
		{"year negative", `{"title":"T","author":"A","year":-5}`, "year must be > 0"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rr := doJSON(t, r, http.MethodPost, "/books", tc.body)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status got %d body=%s", rr.Code, rr.Body.String())
			}
			er := decodeJSON[errResp](t, rr)
			if er.Error != tc.wantError {
				t.Fatalf("error got %q want %q", er.Error, tc.wantError)
			}
		})
	}
}

func TestBooks_BadJSON(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	rr := doJSON(t, r, http.MethodPost, "/books", `{"title":`)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status got %d body=%s", rr.Code, rr.Body.String())
	}
	er := decodeJSON[errResp](t, rr)
	if er.Error != "invalid JSON body" {
		t.Fatalf("error got %q want %q", er.Error, "invalid JSON body")
	}
}

func TestBooks_NotFound(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	// GET unknown id
	rr := doJSON(t, r, http.MethodGet, "/books/9999", ``)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("get status got %d body=%s", rr.Code, rr.Body.String())
	}

	// PUT unknown id
	rr = doJSON(t, r, http.MethodPut, "/books/9999", `{"title":"T","author":"A","year":2000}`)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("put status got %d body=%s", rr.Code, rr.Body.String())
	}

	// DELETE unknown id
	req := httptest.NewRequest(http.MethodDelete, "/books/9999", nil)
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusNotFound {
		t.Fatalf("delete status got %d body=%s", rr2.Code, rr2.Body.String())
	}
}

func TestBooks_InvalidIDParam(t *testing.T) {
	r, db := setupTestRouter(t)
	defer db.Close()

	paths := []string{"/books/abc", "/books/0", "/books/-1"}

	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			rr := doJSON(t, r, http.MethodGet, p, ``)
			if rr.Code != http.StatusBadRequest {
				t.Fatalf("status got %d body=%s", rr.Code, rr.Body.String())
			}
			er := decodeJSON[errResp](t, rr)
			if er.Error != "invalid id" {
				t.Fatalf("error got %q want %q", er.Error, "invalid id")
			}
		})
	}
}
