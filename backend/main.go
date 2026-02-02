package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	// Swagger docs (generated via swag init)
	_ "byfood/backend/docs"
)

// @title byFood Assignment API
// @version 1.0
// @description Books CRUD + URL Processor service
// @BasePath /
func main() {
	// database part
	db, err := OpenDB("file:books.db?_pragma=busy_timeout(5000)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := Migrate(db); err != nil {
		log.Fatal(err)
	}

	store := NewBookStore(db)
	booksAPI := NewBooksAPI(store)

	// router part
	r := chi.NewRouter()

	// Core middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// CORS (in Next.js)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300, // cache preflight for 5 minutes
	}))

	// Swagger
	// http://localhost:8080/swagger/index.html
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Part 2: URL Processor
	r.Post("/process-url", ProcessURLHandler)

	// Part 1: Books CRUD
	r.Route("/books", func(r chi.Router) {
		r.Get("/", booksAPI.GetBooksHandler)
		r.Post("/", booksAPI.CreateBookHandler)
		r.Get("/{id}", booksAPI.GetBookHandler)
		r.Put("/{id}", booksAPI.UpdateBookHandler)
		r.Delete("/{id}", booksAPI.DeleteBookHandler)
	})

	// Start server
	addr := ":8080"
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
