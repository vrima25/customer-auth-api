# Go Auth Service

**Live:** https://go-auth-service-hh8n.onrender.com/
**Health check:** https://go-auth-service-hh8n.onrender.com//health

> Hosted on a free tier. The first request after a period of inactivity may take
> 30–60 seconds while the service wakes up. Subsequent requests are fast.

A lightweight authentication service written in Go, built with the standard `net/http` library — no web framework. It handles customer registration and login using JWT-based stateless authentication and bcrypt password hashing.

---

## Features

- Customer registration with email + password
- Login with JWT token issuance
- Protected endpoint (`/api/v1/customers/me`) guarded by custom auth middleware
- Password hashing with bcrypt — plain-text passwords are never stored
- Health check that verifies the database connection, not just the process
- Layered architecture with dependency injection via interfaces
- PostgreSQL persistence via `database/sql` + `lib/pq`
- Configuration validated at startup — the app refuses to run with a missing secret

---

## Tech Stack

| Layer | Technology |
| --- | --- |
| Language | Go 1.22+ |
| HTTP | Standard library `net/http` |
| Database | PostgreSQL (Neon) |
| Driver | `lib/pq` via `database/sql` |
| Auth | JWT (`golang-jwt/jwt/v5`) |
| Password hashing | bcrypt (`golang.org/x/crypto/bcrypt`) |
| Config | Environment variables, `.env` locally via `joho/godotenv` |
| Hosting | Render (web service) + Neon (database) |

---

## Project Structure

```
.
├── main.go          # entry point: wiring + HTTP server
├── config/          # environment loading and validation
├── model/           # domain structs
├── interfaces/      # contracts between layers
├── repository/      # database access implementation
├── service/         # business logic (hashing, validation, JWT)
├── controller/      # HTTP handlers
├── middleware/      # JWT auth middleware
└── util/            # JWT generate/parse helpers
```

Each layer depends only on the interface of the layer below it, never on its concrete implementation. Changing the database driver does not require touching the service layer.

---

## Design decisions

> These are the choices I made and the trade-offs I accepted.

**Standard library over a framework (Gin, Echo, Fiber)**
I wanted to understand routing, middleware chaining, and request handling directly rather than through a framework's abstractions. The cost is more boilerplate — I write my own JSON helpers and middleware wiring. For a service this size, that cost is low and the understanding is worth it.

**JWT over server-side sessions**
JWTs keep the service stateless, so it can scale horizontally without shared session storage. The trade-off is real: a token cannot be revoked before it expires. I mitigate this with a short 60-minute lifetime. A production system handling sensitive operations would need a refresh-token rotation scheme or a revocation list.

**bcrypt for password hashing**
bcrypt is deliberately slow and salts each password automatically. Slowness is the point — it makes brute-forcing a leaked hash expensive. A fast hash like SHA-256 is wrong for passwords precisely because it is fast.

**Interfaces between layers**
The service layer depends on a `CustomerRepository` interface, not on the concrete PostgreSQL implementation. This keeps business logic testable without a database and makes the storage layer swappable.

**Configuration fails fast at startup**
`config.Load()` returns an error if a required variable is missing. The application refuses to start rather than running with an empty JWT secret and issuing tokens anyone could forge. Failing loudly at boot is safer than failing silently at request time.

**Health check verifies the database**
`/health` pings the database with a 2-second timeout and returns `503` when it is unreachable. A health check that always returns `200` tells the platform everything is fine while the database is down.

---

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL running locally, or a connection string to a hosted instance

### Setup

1. Clone the repository

   ```bash
   git clone https://github.com/vrima25/go-auth-service
   cd go-auth-service
   ```

2. Copy the environment template and fill in your own values

   ```bash
   cp .env.example .env
   ```

3. Create the table

   ```sql
   CREATE TABLE customers (
       id            SERIAL PRIMARY KEY,
       full_name     VARCHAR(100) NOT NULL,
       email         VARCHAR(100) NOT NULL UNIQUE,
       password_hash TEXT NOT NULL,
       created_at    TIMESTAMP DEFAULT now()
   );
   ```

4. Install dependencies and run

   ```bash
   go mod tidy
   go run .
   ```

The server starts on `http://localhost:8080` unless `PORT` is set.

### Environment variables

| Variable | Required | Notes |
| --- | --- | --- |
| `JWT_SECRET` | yes | Generate with `openssl rand -base64 32` |
| `DATABASE_URL` | — | Full connection string. Takes precedence when set. |
| `DB_HOST` | * | Local fallback, default `localhost` |
| `DB_PORT` | * | Local fallback, default `5432` |
| `DB_USER` | * | Required when `DATABASE_URL` is not set |
| `DB_PASSWORD` | * | Local fallback |
| `DB_NAME` | * | Required when `DATABASE_URL` is not set |
| `DB_SSLMODE` | — | Default `disable`. Managed Postgres usually needs `require`. |
| `PORT` | — | Set by the hosting platform. Defaults to `8080`. |

---

## API Reference

### Health check

```
GET /health
```

Verifies the service and its database connection. Returns `503` when the database is unreachable rather than reporting healthy unconditionally.

**Response `200 OK`**

```json
{ "status": "ok", "database": "ok" }
```

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

### Get current customer (protected)

```
GET /api/customers/me
Authorization: Bearer <token>
```

**Response `200 OK`**

```json
{ "email": "jane@example.com" }
```

---

## Security Notes

- Passwords are hashed with bcrypt before storage; plain-text passwords are never persisted.
- JWTs are short-lived (60 minutes) and signed with HMAC-SHA256 using a secret loaded from the environment.
- Login failures return a generic error regardless of whether the email or the password was wrong, so the endpoint cannot be used to enumerate registered accounts.
- The password hash is excluded from all JSON responses.
- Database errors are logged server-side and never returned to the client, since driver messages can expose host and database names.

> This is a public demo instance. Do not register real credentials.

---

## Deployment

Deployed as a web service on Render, with PostgreSQL hosted on Neon. The application reads `PORT` and `DATABASE_URL` from the environment, so no code changes are needed between local and production.

---

## What I'd do next

- Unit tests for the service layer
- Database migrations — the schema is currently applied manually
- Structured logging with request IDs
- Rate limiting on the login endpoint
- Refresh token rotation, so sessions can be revoked before expiry
- Graceful shutdown, so in-flight requests are not dropped on deploy
