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

---

## 📋 Prerequisites

Before running this application, ensure you have:

- **Docker** (>= 20.10)
- **Docker Compose** (>= 2.0) or `docker-compose` (>= 1.29)
- **Go** (>= 1.21) - only for local development
- **Git** - to clone the repository

### Quick Prerequisites Check
```
docker --version
docker-compose --version  # or: docker compose version
go version
```

---

## 🚀 Quick Start

**For the fastest setup, use our automated script:**

### Step 1: Make the setup script executable
```
chmod +x local_setup.sh
```

### Step 2: Run the automated setup
```
./local_setup.sh
```

### Step 3: Access your application
- **App**: http://localhost:8080
- **MySQL Database**: localhost:3307 (external access)

This script will:
- Copy `.env.example` to `.env` if it doesn't exist
- Build and start all services with Docker Compose
- Set up the database automatically
- Get your app running on `http://localhost:8080`

---

## 🚀 Manual Setup (Docker Compose)

If you prefer manual setup:

1. **Copy `.env.example` to `.env` and fill in your DB credentials:**
    ```
    cp .env.example .env
    # Edit .env as needed
    ```

2. **Build and start the app:**
    ```
    docker-compose up --build
    ```

3. **View logs:**
    ```
    docker-compose logs -f
    ```

4. **Stop the app:**
    ```
    docker-compose down
    ```

---

## 🛠️ Local Development (without Docker)

1. **Install Go (>= 1.24) and MySQL locally.**
2. **Create a `.env` file with your DB credentials.**
3. **Run the app:**
    ```
    go mod tidy
    go run main.go
    ```

---

## 🧪 Sample cURL Requests

### Get all books (paginated)
```
curl -X GET "http://localhost:8080/books?page=1&limit=10" \
  -H "Authorization: Bearer this-is-a-secret-token" \
  -H "Content-Type: application/json"
```

### Get a book by ID
```
curl -X GET "http://localhost:8080/books/1" \
  -H "Authorization: Bearer this-is-a-secret-token" \
  -H "Content-Type: application/json"
```

### Create a new book
```
curl -X POST "http://localhost:8080/books" \
  -H "Authorization: Bearer this-is-a-secret-token" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Sample Book Title",
    "author": "Sample Author",
    "genre": "Fiction",
    "published_year": 2023
  }'
```

### Update a book
```
curl -X PUT "http://localhost:8080/books/<BOOK_ID>" \
  -H "Authorization: Bearer this-is-a-secret-token" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Book Title",
    "author": "Updated Author",
    "genre": "Non-Fiction"
  }'
```

### Delete a book
```
curl -X DELETE "http://localhost:8080/books/<BOOK_ID>" \
  -H "Authorization: Bearer this-is-a-secret-token"
```

### Testing without Authorization (should return 401)
```
curl -X GET "http://localhost:8080/books" \
  -H "Content-Type: application/json"
```

---

## ⚙️ Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_USER` | Database username | `root` |
| `DB_PASSWORD` | Database password | `test123test123` |
| `DB_HOST` | Database host (container name) | `mysql` |
| `DB_NAME` | Database name | `digicert` |
| `MYSQL_ROOT_PASSWORD` | MySQL root password | `test123test123` |
| `MYSQL_DATABASE` | MySQL database to create | `digicert` |

---

## ✨ Features

- 📚 Full CRUD operations for books
- 📄 Pagination support for large datasets
- 🐳 Containerized with Docker
- 🔄 Database migrations with Goose
- 🌐 RESTful API design
- 🔒 Environment-based configuration
- 📊 Structured logging
- 🧪 Unit test support

---

## 📝 Notes

- **Environment variables** are managed via `.env` (never commit secrets; use `.env.example` for sharing).
- **Database** is automatically started via Docker Compose (`mysql` service).
- **Tests** are located alongside code in `_test.go` files.

---

## 📚 Pagination Example

To get page 2 with 10 books per page:

GET /books?page=2&limit=10

---


## 🧹 Docker Cleanup (if needed)

If you encounter permission denied issues: Follow these

sudo systemctl stop docker
sudo systemctl start docker

sudo usermod -aG docker $USER
newgrp docker

sudo docker system prune -a --volumes
./local_setup.sh

---

If you encounter issues or want a fresh start:

### Quick Cleanup (Recommended)

# Stop and remove project containers
docker-compose down -v

# Remove project images
docker-compose down --rmi all -v

### Full Docker Cleanup (Use with caution)

# Stop all containers
docker stop $(docker ps -aq)

# Remove all containers
docker rm $(docker ps -aq)

# Remove all images
docker rmi $(docker images -q)

# Clean up everything (containers, images, volumes, networks)
docker system prune -a --volumes

**⚠️ Warning:** Full cleanup will remove ALL Docker data including other projects.

---

## 🧑‍💻 Contributing

1. Fork the repo.
2. Clone and create your `.env` file.
3. Run tests with:
    ```
    go test ./...
    ```
4. Submit a pull request.

---

## 📄 License

MIT
