# Digicert Library App

A scalable Go application for managing a library of books, containerized with Docker and orchestrated using Docker Compose.

---

## 📁 Folder Structure

```
digicert-library-app/
├── internal/
│   ├── database/         # DB connection, queries
│   ├── handlers/
│   │   └── books/        # Books handler logic
│   └── middleware/       # Middlewares (logging, auth, etc.)
├── main.go               # App entry point
├── go.mod, go.sum        # Go modules
├── Dockerfile            # Container build
├── docker-compose.yaml   # Multi-service orchestration
├── .env                  # Environment variables (not committed)
├── .env.example          # Sample env file for setup
└── scripts/              # Setup or utility scripts
```

---

## 🚦 API Routes

| Method | Route             | Description                  |
|--------|-------------------|------------------------------|
| GET    | `/`               | Welcome message              |
| GET    | `/books`          | List books (supports pagination) |
| GET    | `/books/{id}`     | Get book by ID               |
| POST   | `/books`          | Create a new book            |
| PUT    | `/books/{id}`     | Update a book                |
| DELETE | `/books/{id}`     | Delete a book                |

---
## 🚀 Quick Start

**For the fastest setup, use our automated script:**

```bash
./scripts/local_setup.sh
```

This script will:
- Copy `.env.example` to `.env` if it doesn't exist
- Build and start all services with Docker Compose
- Set up the database automatically
- Get your app running on `http://localhost:8080`

---

## 🚀 Manual Setup (Docker Compose)

If you prefer manual setup:

1. **Copy `.env.example` to `.env` and fill in your DB credentials:**
    ```bash
    cp .env.example .env
    # Edit .env as needed
    ```

2. **Build and start the app:**
    ```bash
    docker-compose up --build
    ```

3. **View logs:**
    ```bash
    docker-compose logs -f
    ```

4. **Stop the app:**
    ```bash
    docker-compose down
    ```

---

## 🛠️ Local Development (without Docker)

1. **Install Go (>= 1.24) and MySQL locally.**
2. **Create a `.env` file with your DB credentials.**
3. **Run the app:**
    ```bash
    go mod tidy
    go run main.go
    ```

---

## 🧪 Sample cURL Requests

### Get all books (paginated)
```bash
curl -X GET "http://localhost:8080/books?page=1&limit=10"
```

### Get a book by ID
```bash
curl -X GET "http://localhost:8080/books/<BOOK_ID>"
```

### Create a new book
```bash
curl -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{"title": "Sample Book Title"}'
```

### Update a book
```bash
curl -X PUT "http://localhost:8080/books/<BOOK_ID>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Book Title"}'
```

### Delete a book
```bash
curl -X DELETE "http://localhost:8080/books/<BOOK_ID>"
```
```

---

## 📝 Notes

- **Environment variables** are managed via `.env` (never commit secrets; use `.env.example` for sharing).
- **Database** is automatically started via Docker Compose (`mysql` service).
- **Tests** are located alongside code in `_test.go` files.
- **Scripts** for setup are in the `scripts/` folder.

---

## 📚 Pagination Example

To get page 2 with 10 books per page:
```
GET /books?page=2&limit=10
```

---

## 🧑‍💻 Contributing

1. Fork the repo.
2. Clone and create your `.env` file.
3. Run tests with:
    ```bash
    go test ./...
    ```
4. Submit a pull request.

---

## 📄 License

MIT
