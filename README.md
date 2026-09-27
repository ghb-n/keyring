# Backend

A RESTful API built with Go, Gin, and PostgreSQL.

## Features

- User registration
- User login with JWT authentication
- Password hashing with bcrypt
- PostgreSQL connection pool (pgx)
- Database migrations (Goose)
- Graceful shutdown

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **DB Driver**: pgx/v5 + pgxpool
- **Migrations**: Goose
- **Auth**: JWT (HS256) + bcrypt
- **Config**: godotenv

## Project Structure

```text
.
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── auth/                 # JWT generation and validation
│   ├── config/               # Environment config loading
│   ├── db/                   # Database connection pool
│   ├── handlers/             # HTTP handlers
│   ├── middleware/           # Gin middleware (auth, etc)
│   ├── models/               # Domain models
│   ├── repository/           # Data access layer
│   ├── routes/               # Route registration
│   └── services/             # Business logic
├── migrations/               # SQL migrations (Goose)
├── .env.example              # Environment template
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Goose

### Setup

1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd Backend
   ```

1. Copy the environment template:

   ```bash
   cp .env.example .env
   ```

2. Fill in the `.env` values:

   ```bash
   DB_PASSWORD=your-password
   JWT_SECRET=your-secret
   ```

   To generate a JWT secret:

   ```bash
   openssl rand -base64 32 | tr '+/' '-_' | tr -d '='
   ```

3. Start PostgreSQL (Docker):

   ```bash
   docker run --name postgres-dev \
     -e POSTGRES_PASSWORD=your-password \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_DB=gin_study \
     -p 5432:5432 \
     -d postgres:16
   ```

4. Install Goose:

   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

5. Run migrations:

   ```bash
   goose -dir migrations postgres "postgres://postgres:your-password@localhost:5432/gin_study?sslmode=disable" up
   ```

6. Install dependencies:

   ```bash
   go mod download
   ```

7. Run the server:

   ```bash
   go run ./cmd/server
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

### Public

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/register` | Create a new user |
| `POST` | `/login` | Authenticate and receive a JWT |

### Protected

Requires `Authorization: Bearer <token>` header.

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/api/me` | Get current user profile |
| `PUT` | `/api/me` | Update current user profile |
| `DELETE` | `/api/me` | Delete current user |

## Examples

### Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Grigoriy","email":"g@example.com","password":"secret123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"g@example.com","password":"secret123"}'
```

### Get current user

```bash
curl http://localhost:8080/api/me \
  -H "Authorization: Bearer <token>"
```

## License
[MIT](LICENSE)
