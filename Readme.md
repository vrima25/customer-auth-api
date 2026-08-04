# Customer Auth API

A lightweight authentication service written in Go, built with the standard `net/http` library (no web framework). It handles customer registration and login using JWT-based stateless authentication and bcrypt password hashing.

## Features

- Customer registration with email + password
- Login with JWT token issuance
- Protected endpoint example (`/me`) guarded by custom auth middleware
- Password hashing with bcrypt (no plain-text storage)
- Clean, layered architecture with dependency injection via interfaces
- PostgreSQL persistence via `database/sql` + `lib/pq`

## Tech Stack

| Layer            | Technology                            |
| ---------------- | ------------------------------------- |
| Language         | Go                                    |
| HTTP             | Standard library `net/http`           |
| Database         | PostgreSQL                            |
| Auth             | JWT (`golang-jwt/jwt/v5`)             |
| Password hashing | bcrypt (`golang.org/x/crypto/bcrypt`) |
| Config           | `.env` via `joho/godotenv`            |

## Project Structure

```
.
├── main.go              # entry point: wiring + HTTP server
├── config/               # environment variable loading
├── model/                # domain structs
├── interfaces/            # contracts between layers
├── repository/            # database access implementation
├── service/               # business logic (hashing, validation, JWT)
├── controller/             # HTTP handlers
├── middleware/             # JWT auth middleware
└── util/                  # JWT generate/parse helpers
```

Each layer depends only on the interface of the layer below it, not on its concrete implementation — this keeps the codebase testable and swappable (e.g. changing the database driver doesn't require touching the service layer).

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL running locally (or accessible via network)

### Setup

1. Clone the repository

   ```bash
   git clone <your-repo-url>
   cd customer-auth-api
   ```

2. Copy the environment template and fill in your own values

   ```bash
   cp .env.example .env
   ```

3. Create the database table

   ```sql
   CREATE TABLE customers (
       id            SERIAL PRIMARY KEY,
       email         VARCHAR(255) UNIQUE NOT NULL,
       password_hash TEXT NOT NULL,
       full_name     VARCHAR(150) NOT NULL,
       created_at    TIMESTAMP DEFAULT now()
   );
   ```

4. Install dependencies

   ```bash
   go mod tidy
   ```

5. Run the server
   ```bash
   go run main.go
   ```

The server starts on `http://localhost:8080`.

## API Reference

### Register

```
POST /api/customers/register
```

**Request body**

```json
{
  "email": "jane@example.com",
  "password": "securepassword123",
  "full_name": "Jane Doe"
}
```

**Response `201 Created`**

```json
{
  "id": 1,
  "email": "jane@example.com",
  "full_name": "Jane Doe"
}
```

### Login

```
POST /api/customers/login
```

**Request body**

```json
{
  "email": "jane@example.com",
  "password": "securepassword123"
}
```

**Response `200 OK`**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "email": "jane@example.com",
  "full_name": "Jane Doe"
}
```

### Get current user (protected)

```
GET /api/customers/me
Authorization: Bearer <token>
```

**Response `200 OK`**

```json
{
  "email": "jane@example.com"
}
```

## Security Notes

- Passwords are hashed with bcrypt before storage; plain-text passwords are never persisted.
- JWTs are short-lived (60 minutes) and signed with HMAC (HS256) using a secret loaded from the environment.
- Login failures return a generic error message regardless of whether the email or password was incorrect, to avoid leaking which field was wrong.
