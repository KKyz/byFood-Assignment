# byFood Assignment – Full Stack CRUD App

This project is a full-stack web application built as part of the byFood assignment.
It implements a RESTful backend API with persistent storage and a frontend dashboard
for managing books.

The application supports full CRUD operations, URL processing, API documentation
via Swagger, and a modern frontend built with Next.js.

---

## Tech Stack

### Backend
- Go
- chi (HTTP router)
- SQLite (persistent storage)
- Swagger (OpenAPI documentation)

### Frontend
- Next.js (App Router)
- React
- TypeScript
- Tailwind CSS
- Context API (global state management)

---

## Features

### Backend
- CRUD API for books
- SQLite persistence
- URL processing endpoint (`/process-url`)
- Input validation and error handling
- Swagger UI for API documentation
- Unit and integration tests

### Frontend
- Dashboard listing all books
- Add / edit / delete books via modal form
- Client-side form validation
- Dynamic routing for book detail pages (`/books/[id]`)
- Global state via Context API
- User-friendly error handling

---

## Backend Setup

```bash
cd backend
go run .
```

The backend runs on:
- **API**: http://localhost:8080 (locally)
- **Swagger UI**: http://localhost:8080/swagger/index.html (locally)

---

## Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

The frontend runs on:
- **App**: http://localhost:3000 (locally)

---
## Running Tests

Run all backend tests:
```bash
cd backend
go test ./...
```

Expected output:
```bash
ok  	byfood/backend
?   	byfood/backend/docs	[no test files]
```
---

## API Documentation

### Books API

#### GET /books
Returns all books.

**Request:**
```bash
curl http://localhost:8080/books
```

**Response:**
```json
[
  {
    "id": 1,
    "title": "Dune",
    "author": "Frank Herbert",
    "year": 1965
  }
]
```

---

#### POST /books
Create a new book.

**Request:**
```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune","author":"Frank Herbert","year":1965}'
```

**Response:**
```json
{
  "id": 1,
  "title": "Dune",
  "author": "Frank Herbert",
  "year": 1965
}
```

---

#### GET /books/{id}
Fetch a single book by ID.

**Request:**
```bash
curl http://localhost:8080/books/1
```

**Error cases:**
- `400 Bad Request` – invalid ID
- `404 Not Found` – book does not exist

---

#### PUT /books/{id}
Update a book.

**Request:**
```bash
curl -X PUT http://localhost:8080/books/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Dune Messiah","author":"Frank Herbert","year":1969}'
```

---

#### DELETE /books/{id}
Delete a book.

**Request:**
```bash
curl -X DELETE http://localhost:8080/books/1
```

**Response:**
```
204 No Content
```

---

### URL Processing API

#### POST /process-url
Processes a URL using one of the supported operations:
- `canonical`
- `redirection`
- `all`

**Request:**
```bash
curl -X POST http://localhost:8080/process-url \
  -H "Content-Type: application/json" \
  -d '{"url":"https://BYFOOD.com/food-EXPeriences?query=abc","operation":"all"}'
```

**Response:**
```json
{
  "processed_url": "https://www.byfood.com/food-experiences"
}
```

**Error cases:**
- Missing required fields → `400 Bad Request`
- Invalid operation → `400 Bad Request`
- Invalid URL → `400 Bad Request`

---

## Project Structure

```bash
byfood-assignment/
├── backend/
│   ├── main.go              # Server entry point (routes, middleware, CORS)
│   ├── db.go                # SQLite connection + migrations
│   ├── books_store.go       # Data access layer (SQL queries)
│   ├── books_handlers.go    # HTTP handlers for /books endpoints
│   ├── url_processor.go     # /process-url endpoint logic
│   ├── *_test.go            # Backend unit & integration tests
│   ├── docs/                # Auto-generated Swagger files
│   ├── go.mod / go.sum
│   └── books.db             # Runtime DB (gitignored)
│
├── frontend/
│   ├── app/
│   │   ├── layout.tsx       # Root layout + Context provider
│   │   ├── page.tsx         # Dashboard (book list)
│   │   └── books/[id]/      # Dynamic route for book detail pages
│   │       └── page.tsx
│   ├── src/
│   │   ├── components/      # Reusable UI components (modal, errors)
│   │   ├── context/         # React Context (global book state)
│   │   └── lib/
│   │       └── api.ts       # Centralized API client
│   ├── .env.local           # Frontend environment variables
│   └── package.json
│
├── screenshots/             # Screenshots for README
└── README.md
```
---

---


## Screenshots

### Dashboard
![Dashboard](screenshots/Dashboard.png)

### Book Detail Page
![Book Detail Page](screenshots/BookDetailPage.png)

### Add Book Modal
![Add Book Modal](screenshots/AddModal.png)

### Edit Book Modal
![Edit Book Modal](screenshots/EditModal.png)

### Swagger UI
![Swagger UI](screenshots/SwaggerUI.png)
