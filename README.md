# 📓 Keeper API

A clean, production-ready REST API for managing users and notes — built with Go, Chi, and PostgreSQL.

---

## ✨ Features

- 🏗️ **Layered architecture** — handlers → services → repositories, clearly separated
- 📄 **Pagination** — all list endpoints support `page` & `limit` query params
- 🛡️ **Structured error handling** — no raw DB errors ever reach the client
- ✅ **Field-level validation errors** — per-field messages on invalid input
- 🔄 **Graceful shutdown** — handles `SIGINT`/`SIGTERM` with a 10s drain window
- 🗄️ **PostgreSQL** with connection pooling via `pgx/v5`

---

## 🗂️ Project Structure

```
keeper/
├── cmd/
│   └── api/
│       └── main.go            # Entry point, server setup, routing
├── internal/
│   ├── configs/
│   │   └── config.go          # Env-based config loading with defaults
│   ├── database/
│   │   └── postgres.go        # pgx connection pool setup
│   ├── di/
│   │   └── container.go       # Dependency injection container
│   ├── errors/
│   │   └── errors.go          # Sentinel errors, DB/validation error mapping
│   ├── handlers/
│   │   ├── note_handler.go    # HTTP handlers for /notes
│   │   └── user_handler.go    # HTTP handlers for /users
│   ├── middlewares/           # Custom middleware (WIP)
│   ├── models/
│   │   ├── note.go            # Note model, input types, validation
│   │   └── user.go            # User model, input types, validation
│   ├── repositories/
│   │   ├── note_repo.go       # Note DB queries (PostgreSQL)
│   │   └── user_repo.go       # User DB queries (PostgreSQL)
│   ├── response/
│   │   └── response.go        # Unified JSON response helpers
│   └── services/
│       ├── note_service.go    # Note business logic
│       └── user_service.go    # User business logic
├── migrations/                # SQL migrations (go-migrate)
├── pkg/
│   └── log/
│       └── zap.go             # Zap logger setup
├── .env.example               # Environment variable template
├── Makefile                   # Development & deployment commands
└── go.mod
```

---

## ⚙️ Configuration

Copy `.env.example` to `.env` and fill in the values:

```bash
cp .env.example .env
```

| Variable        | Default         | Description                        |
|-----------------|-----------------|------------------------------------|
| `APP_PORT`      | `3000`          | Port the server listens on         |
| `DB_URL`        | *(required)*    | PostgreSQL connection string        |
| `ENVIRONMENT`   | `development`   | Runtime environment label          |
| `READ_TIMEOUT`  | `15`            | HTTP read timeout in seconds       |
| `WRITE_TIMEOUT` | `30`            | HTTP write timeout in seconds      |
| `IDLE_TIMEOUT`  | `30`            | HTTP idle timeout in seconds       |

**Example `DB_URL`:**
```
DB_URL=postgres://user:password@localhost:5432/keeper?sslmode=disable
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI

### Install & Run

```bash
# 1. Clone the repo
git clone https://github.com/kevalsabhani/keeper.git
cd keeper

# 2. Install dependencies
make deps

# 3. Set up environment
cp .env.example .env
# Edit .env with your DB_URL and other values

# 4. Run database migrations
make migrate-up

# 5. Start the server
make run
```

---

## 🛠️ Makefile Commands

```bash
make help           # Show all available commands
make build          # Compile the API binary to ./bin/api
make run            # Build and start the server
make dev            # Run without building (go run)
make test           # Run all tests with race detector
make test-coverage  # Run tests and open coverage report in browser
make lint           # Run golangci-lint
make fmt            # Format code with goimports
make clean          # Remove build artifacts and coverage files
make migrate-up     # Apply all pending migrations
make migrate-down   # Roll back all migrations
make docker-build   # Build Docker image
make docker-up      # Start services via docker-compose
make docker-down    # Stop docker-compose services
make deps           # Run go mod tidy && go mod verify
```

---

## 🌐 API Reference

**Base URL:** `http://localhost:3000/api/v1`

### System

| Method | Endpoint    | Description         |
|--------|-------------|---------------------|
| `GET`  | `/`         | Welcome message     |
| `GET`  | `/health`   | Health check        |

---

### 👤 Users

| Method   | Endpoint       | Description              |
|----------|----------------|--------------------------|
| `POST`   | `/users`       | Create a user            |
| `GET`    | `/users`       | List users (paginated)   |
| `GET`    | `/users/{id}`  | Get user by ID           |
| `PATCH`  | `/users/{id}`  | Update a user            |
| `DELETE` | `/users/{id}`  | Delete a user            |

#### Create User
```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "Alice",
  "email": "alice@example.com"
}
```

#### List Users (paginated)
```bash
GET /api/v1/users?page=1&limit=10
```

---

### 📝 Notes

| Method   | Endpoint       | Description              |
|----------|----------------|--------------------------|
| `POST`   | `/notes`       | Create a note            |
| `GET`    | `/notes`       | List notes (paginated)   |
| `GET`    | `/notes/{id}`  | Get note by ID           |
| `PATCH`  | `/notes/{id}`  | Update a note            |
| `DELETE` | `/notes/{id}`  | Delete a note            |

#### Create Note
```bash
POST /api/v1/notes
Content-Type: application/json

{
  "user_id": 1,
  "title": "My first note",
  "content": "Hello, Keeper!"
}
```

---

## 📦 Response Format

All endpoints return a consistent JSON envelope:

### Success
```json
{
  "success": true,
  "data": [ { "id": 1, "name": "Alice", "email": "alice@example.com", "..." } ],
  "meta": {
    "current_page": 1,
    "total_pages": 4,
    "total_count": 37
  }
}
```

> Single-resource endpoints (e.g. `GET /users/{id}`) also return `data` as an array with one item.

### Error
```json
{
  "success": false,
  "error": {
    "code": "INVALID_INPUT",
    "message": "invalid input",
    "fields": [
      { "field": "Title",   "message": "Title is required" },
      { "field": "Content", "message": "Content must be at least 3 characters" }
    ]
  }
}
```

> `fields` is only present for validation errors (`400`). Other errors return only `code` and `message`.

### Error Codes

| Code                    | HTTP Status | When                                    |
|-------------------------|-------------|-----------------------------------------|
| `INVALID_INPUT`         | `400`       | Validation failure or bad request body  |
| `NOT_FOUND`             | `404`       | Resource does not exist                 |
| `CONFLICT`              | `409`       | Duplicate unique field (e.g. email)     |
| `INTERNAL_SERVER_ERROR` | `500`       | Unexpected server-side error            |

---

## 🧱 Tech Stack

| Layer       | Technology                                                                 |
|-------------|----------------------------------------------------------------------------|
| Language    | [Go](https://go.dev)                                                       |
| Router      | [chi v5](https://github.com/go-chi/chi)                                    |
| Database    | [PostgreSQL](https://www.postgresql.org) via [pgx v5](https://github.com/jackc/pgx) |
| Validation  | [go-playground/validator v10](https://github.com/go-playground/validator)  |
| Logging     | [Uber Zap](https://github.com/uber-go/zap)                                 |
| Migrations  | [golang-migrate](https://github.com/golang-migrate/migrate)                |

---

## 🗃️ Database Schema

```sql
-- Users
CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);

-- Notes
CREATE TABLE notes (
    id         INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    title      VARCHAR(100) NOT NULL,
    content    VARCHAR(100) NOT NULL,
    user_id    INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_notes_user_id ON notes(user_id);
```

---

## 📜 License

MIT
