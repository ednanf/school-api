# School Management API

A RESTful backend for the administrative side of a school — students, classes, subjects, teachers, staff and user accounts — built with Go, Chi and MariaDB, secured with JWT and Argon2id.

[![Go](https://img.shields.io/badge/go-1.27-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Chi](https://img.shields.io/badge/router-Chi%20v5-00ADD8?style=flat-square)](https://github.com/go-chi/chi)
[![MariaDB](https://img.shields.io/badge/database-MariaDB-00ADD8?style=flat-square&logo=mariadb&logoColor=white)](https://mariadb.org/)
[![OpenAPI](https://img.shields.io/badge/spec-OpenAPI%203.0-6BA539?style=flat-square&logo=openapi&logoColor=white)](openapi.yaml)
[![License](https://img.shields.io/badge/license-MIT-6BA539?style=flat-square)](LICENSE)

---

## About this project

Most school software is either a bloated all-in-one suite or a spreadsheet with an API bolted on. This project is a focused backend for the operational core of a school: who the students are, which class they belong to, what is being taught, who teaches what, and who works in the building.

It is a **backend-only** service. There is no UI here — the API is designed to sit behind a separate frontend, and everything is shaped around that: JSON everywhere, a consistent response envelope, pagination on every collection endpoint, and CORS configured for a browser client on a different origin.

The interesting part is not the CRUD. It is the plumbing around it:

- **A real authentication lifecycle.** Accounts are created inactive, activated through a one-time token delivered by e-mail, and can recover a forgotten password through a flow that is deliberately blind to whether the address exists.
- **Passwords hashed with Argon2id** using OWASP-aligned parameters, stored in standard PHC format, and compared in constant time.
- **Rate limiting on the sensitive public routes** (login, password reset, account activation) with a per-IP token bucket, so brute-force attempts and mail-bombing a mailbox both hit a wall.
- **Role-based access control** with a middleware that gates whole route groups, not individual handlers.
- **Transactional bulk operations** — importing 2,000 students is one transaction, so a failure halfway through rolls the whole batch back instead of leaving a half-imported class.

The codebase is organised in four layers with dependencies pointing inwards only, and the domain layer owns the interfaces that the outer layers implement. That means the HTTP handlers and the SQL repositories are both swappable without touching business logic, and every concrete implementation is asserted at compile time the moment it is constructed.

API documentation can be found [here](https://docs.frizzera.dev/api/school-api).

> **Project status:** functionally complete and stable. The code is currently in a readability pass — naming, structure and duplication are being tightened, but the HTTP contract and behavior are unchanged.

---

## Table of contents

- [About this project](#about-this-project)
- [Features](#features)
- [Tech stack](#tech-stack)
- [Architecture](#architecture)
- [Project structure](#project-structure)
- [Getting started](#getting-started)
- [API overview](#api-overview)
- [Conventions](#conventions)
- [Security](#security)
- [Testing](#testing)
- [License](#license)

---

## Features

### Authentication & access control

- **JWT authentication** — HS256 tokens carrying `user_id`, `email` and `role` claims, verified on every protected route.
- **Login** — e-mail + password, with `last_login_at` maintained on every successful sign-in.
- **Account activation** — newly created accounts start inactive and unlock through a single-use token sent by e-mail.
- **Password recovery** — forgot-password → reset token → new password, with token expiry enforced server-side.
- **Change password** — requires the current password, and clears any outstanding reset token.
- **Role-based access control** — three roles today (`super_admin`, `admin`, `staff`), enforced by a `RequireRole` middleware that protects entire route groups. A broader role matrix is documented in [`docs/roles.md`](https://github.com/ednanf/school-api/blob/master/docs/roles.md).
- **Enumeration-safe recovery** — forgot-password always returns the same response whether or not the account exists.
- **Per-IP rate limiting** — separate buckets for login and for the recovery/activation endpoints.

### Core resources

| Resource                | What it covers                                                                                           |
| ----------------------- | -------------------------------------------------------------------------------------------------------- |
| **Students**            | Enrollment records tied to a class, with bulk import, bulk edit, bulk delete and bulk class reassignment |
| **Classes**             | Grade + letter cohorts (Grade 1–9, sections A–D) and a roster view of the students inside them           |
| **Subjects**            | The subject catalogue offered across grades                                                              |
| **Teachers**            | Teaching staff records                                                                                   |
| **Teacher assignments** | Which teacher teaches which subject in which class — the join that makes a timetable possible            |
| **Departments**         | Organizational units such as Science or Student Services                                                 |
| **Staff positions**     | Job titles, each owned by a department                                                                   |
| **Staff**               | Employees, each attached to a position, hire date included                                               |
| **Users**               | Login accounts, their role and their activation state (super-admin only)                                 |

### API ergonomics

- **Consistent response envelope** for every endpoint — success, error and paginated responses all share one shape.
- **Pagination** on all collection endpoints, with `page`/`limit` in, and `page`, `limit`, `count`, `total_items`, `total_pages` out.
- **Structured validation errors** — field-level `422` responses that tell the client exactly which field failed and why.
- **Populated reads** — a student response embeds its class, an assignment embeds teacher + class + subject, a staff member embeds position and department. One round trip, no N+1 in the client.
- **PATCH semantics** — pointer-based patch DTOs, so an omitted field is genuinely absent and a `false` is not mistaken for "not provided".
- **Input normalization** — names and e-mails are trimmed, lower-cased and title-cased consistently before they reach the database.
- **OpenAPI 3.0 specification** covering the full contract, kept alongside the code in [`openapi.yaml`](https://github.com/ednanf/school-api/blob/master/openapi.yaml).

### Operations

- **Idempotent admin bootstrap** — the first `super_admin` is seeded from environment variables on boot, and the seed skips itself if one already exists.
- **Swappable e-mail provider** — Mailpit over plain SMTP for local development, Resend for production, behind one `domain.Mailer` interface.
- **Docker Compose stack** for MariaDB and Mailpit so a contributor can be running in minutes.
- **Seeded schema and reference data** — the full DDL plus starter subjects, classes, departments, positions and staff in [`init.sql`](https://github.com/ednanf/school-api/blob/master/init.sql).

---

## Tech stack

**Language:** Go 1.27 · **Database:** MariaDB / MySQL · **Style:** REST + JSON over HTTP/1.1

| Dependency                                                                             | Role in the project                                                                                                   |
| -------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| [`github.com/go-chi/chi/v5`](https://github.com/go-chi/chi)                            | HTTP router and middleware stack. Chosen for composable route groups and its small, dependency-light footprint.       |
| [`github.com/go-chi/cors`](https://github.com/go-chi/cors)                             | CORS handling for browser clients on another origin.                                                                  |
| [`github.com/go-playground/validator/v10`](https://github.com/go-playground/validator) | Struct-tag driven request validation, including `dive` rules for validating slices of DTOs.                           |
| [`github.com/golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)                    | JWT signing and verification (HS256), with the signing method pinned to prevent algorithm downgrade.                  |
| [`github.com/jmoiron/sqlx`](https://github.com/jmoiron/sqlx)                           | SQL ergonomics over `database/sql`: named parameters, `IN` clause expansion, query rebinding and transaction helpers. |
| [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql)             | MySQL / MariaDB driver, registered with `database/sql`.                                                               |
| [`golang.org/x/crypto`](https://pkg.go.dev/golang.org/x/crypto)                        | Argon2id password hashing.                                                                                            |
| [`golang.org/x/time`](https://pkg.go.dev/golang.org/x/time)                            | Token-bucket rate limiting used by the per-IP limiter.                                                                |
| [`golang.org/x/text`](https://pkg.go.dev/golang.org/x/text)                            | Locale-aware title casing, so names are stored consistently.                                                          |
| [`github.com/resend/resend-go/v4`](https://github.com/resend/resend-go)                | Transactional e-mail delivery in production.                                                                          |
| [`github.com/joho/godotenv`](https://github.com/joho/godotenv)                         | Loads `.env` for local development.                                                                                   |

Standard library highlights: `crypto/rand` for token generation, `crypto/subtle` for constant-time hash comparison, `net/smtp` for the local Mailpit client, and `net/http` for the server itself.

---

## Architecture

The service is a four-layer monolith. Each layer depends only on the one below it, and the domain layer defines the interfaces the outer layers implement — so the transport and persistence details are swappable without touching the business rules.

```
                         HTTP request
                              │
                              ▼
        ┌────────────────────────────────────────────┐
        │  internal/transport/http                   │
        │  Chi router · middleware · handlers · DTOs │
        │  decode → validate → call service → respond│
        └────────────────────┬───────────────────────┘
                             │  domain.*Service
                             ▼
        ┌────────────────────────────────────────────┐
        │  internal/service                          │
        │  business rules · normalization · merging  │
        │  pagination · timestamps · orchestration   │
        └────────────────────┬───────────────────────┘
                             │  domain.*Repository
                             ▼
        ┌────────────────────────────────────────────┐
        │  internal/repository                       │
        │  SQL · named params · transactions         │
        └────────────────────┬───────────────────────┘
                             ▼
                     MariaDB / MySQL
```

Alongside the layers sits a set of small, self-contained packages under `internal/pkg` — password hashing, JWT, e-mail delivery, rate limiting, secure token generation, the admin seeder and a debug logger. Each one is a focused unit with a narrow interface defined in `internal/domain`, which is what keeps the service layer honest about what it actually depends on.

**Key design decisions**

- **Interfaces live in the domain layer.** `domain.StudentRepository` and `domain.StudentService` are declared next to the `Student` entity, not next to their implementations. Both concrete types are returned as interfaces from their constructors, so the compiler checks the contract at the moment of construction.
- **The handler does three things only:** decode, validate, delegate. Every decision about _what should happen_ belongs to the service.
- **The service owns normalization and timestamps.** Repositories receive entities that are already trimmed, cased and stamped, which keeps SQL dumb and makes the behavior testable without a database.
- **Errors are domain-level.** `domain.ErrNotFound`, `ErrInvalidCredentials` and friends are matched with `errors.Is` in the transport layer, and SQL error codes never leak past the repository boundary.
- **The response envelope is centralized.** `sendSuccess`, `sendPaginated` and `sendError` are the only three ways this API writes a body.

---

## Project structure

```
school-api/
├── cmd/
│   └── api/
│       └── main.go          # entrypoint: config, DI wiring, routes, server
├── internal/
│   ├── config/
│   │   └── cors.go          # CORS policy
│   ├── domain/              # entities, DTOs, errors, interfaces
│   │   ├── auth.go          # token + auth service contracts, auth DTOs
│   │   ├── student.go       # student entity, patch & bulk DTOs
│   │   ├── class.go         # class entity and patch DTO
│   │   ├── teacher.go       # teacher entity and patch DTO
│   │   ├── teacher_assignment.go
│   │   ├── subject.go
│   │   ├── department.go
│   │   ├── staff.go
│   │   ├── staff_position.go
│   │   ├── user.go          # user entity, role constants, sanitized response
│   │   ├── password.go      # PasswordHasher contract
│   │   ├── mailer.go        # Mailer contract
│   │   └── errors.go        # sentinel errors shared across layers
│   ├── pkg/
│   │   ├── hasher/          # Argon2id implementation
│   │   ├── token/           # JWT generation and validation
│   │   ├── mailer/          # Mailpit (SMTP) and Resend providers
│   │   ├── ratelimit/       # per-IP token bucket + tests
│   │   ├── security/        # cryptographically secure token generation
│   │   ├── seeder/          # idempotent initial admin bootstrap
│   │   └── loggers/         # small debug logging helper
│   ├── repository/          # MariaDB persistence, one file per resource
│   ├── service/             # business logic, one file per resource
│   └── transport/
│       └── http/            # handlers, middleware, response envelope
├── docs/
│   ├── roles.md             # intended role model and permissions
│   ├── students_dummy.json  # sample bulk import payloads
│   └── subjects_dummy.json  # sample subject payloads
├── init.sql                 # schema + reference data
├── openapi.yaml             # OpenAPI 3.0 specification
├── docker-compose.yml       # MariaDB + Mailpit
└── go.mod
```

---

## Getting started

### Prerequisites

- Go 1.27 or newer
- Docker and Docker Compose (for MariaDB and Mailpit)

### 1. Clone and configure

```bash
git clone https://github.com/ednanf/school-api.git
cd school-api
```

Create a `.env` file in the project root — the table below lists every key the API reads:

```bash
touch .env
```

| Variable                          | Purpose                                                          | Example                                  |
| --------------------------------- | ---------------------------------------------------------------- | ---------------------------------------- |
| `API_PORT`                        | Port the API listens on                                          | `:8000`                                  |
| `ENVIRONMENT`                     | `development` or `production`; selects the default mail provider | `development`                            |
| `DB_USER` / `DB_PASSWORD`         | Database credentials                                             | `admin` / `letmeinadmin`                 |
| `DB_HOST` / `DB_PORT` / `DB_NAME` | Database location and schema                                     | `127.0.0.1` / `3306` / `school_db`       |
| `JWT_SECRET`                      | Signing key for access tokens — **required**                     | _(generate your own)_                    |
| `INITIAL_ADMIN_USERNAME`          | Username for the seeded `super_admin`                            | `super_admin`                            |
| `INITIAL_ADMIN_EMAIL`             | E-mail for the seeded `super_admin`                              | `admin@school.edu`                       |
| `INITIAL_ADMIN_PASSWORD`          | Password for the seeded `super_admin`                            | _(a strong one)_                         |
| `EMAIL_PROVIDER`                  | `mailpit` for local, `resend` for production                     | `mailpit`                                |
| `MAIL_FROM`                       | Sender address for outgoing e-mail                               | `no-reply@example.com`                   |
| `MAILPIT_HOST` / `MAILPIT_PORT`   | Local SMTP server for Mailpit                                    | `localhost` / `1025`                     |
| `MAILPIT_RESET_URL`               | Frontend URL the reset/activation link points at                 | `http://localhost:5173/activate?token`   |
| `RESEND_API_KEY`                  | API key for Resend                                               | `re_...`                                 |
| `RESEND_RESET_URL`                | Frontend URL used when sending via Resend                        | `https://app.example.com/activate?token` |

### 2. Start the dependencies

```bash
docker compose up -d
```

This brings up MariaDB on `3306` and Mailpit on `1025` (SMTP) with a web inbox at [http://localhost:8025](http://localhost:8025).

### 3. Load the schema and reference data

```bash
docker compose exec -T mariadb mariadb -u admin -pletmeinadmin < init.sql
```

### 4. Run the API

```bash
go run ./cmd/api
```

The server starts on `:8000` and seeds the initial `super_admin` account from the environment variables above — if one already exists, the seed quietly skips itself.

### 5. Log in and make a request

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@school.edu","password":"SuperSecurePassword123!"}'
```

```json
{
  "status": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 100,
      "username": "super_admin",
      "email": "admin@school.edu",
      "role": "super_admin",
      "is_active": true
    }
  }
}
```

Use the token on any protected endpoint:

```bash
curl http://localhost:8000/api/v1/students?page=1&limit=10 \
  -H "Authorization: Bearer <your_token>"
```

---

## API overview

All endpoints are prefixed with `/api/v1` and return JSON.

### Auth — `/api/v1/auth`

| Method | Path                    | Description                                  | Auth                  |
| ------ | ----------------------- | -------------------------------------------- | --------------------- |
| `POST` | `/auth/login`           | Authenticate and receive a JWT               | Public · rate limited |
| `POST` | `/auth/forgot-password` | Request a password reset link                | Public · rate limited |
| `POST` | `/auth/reset-password`  | Consume a reset token and set a new password | Public · rate limited |
| `POST` | `/auth/activate`        | Activate a newly created account             | Public · rate limited |
| `POST` | `/auth/change-password` | Change the password of the signed-in user    | Bearer token          |

### Resources

| Resource                | Collection                    | Operations                                                                   |
| ----------------------- | ----------------------------- | ---------------------------------------------------------------------------- |
| **Users**               | `/api/v1/users`               | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}` — _super-admin only_   |
| **Students**            | `/api/v1/students`            | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |
| **Students (bulk)**     | `/api/v1/students`            | `POST /bulk` · `PATCH /bulk` · `DELETE /bulk` · `PATCH /class`               |
| **Classes**             | `/api/v1/classes`             | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}` · `GET /{id}/students` |
| **Subjects**            | `/api/v1/subjects`            | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |
| **Teachers**            | `/api/v1/teachers`            | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |
| **Teacher assignments** | `/api/v1/teacher_assignments` | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |
| **Departments**         | `/api/v1/departments`         | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |
| **Staff positions**     | `/api/v1/staff_positions`     | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |
| **Staff**               | `/api/v1/staff`               | `GET` `POST` `GET /{id}` `PATCH /{id}` `DELETE /{id}`                        |

Every endpoint except the public auth routes requires `Authorization: Bearer <token>`. The full request and response schemas live in [`openapi.yaml`](openapi.yaml) — import it into Swagger UI, Postman or Insomnia to explore the API interactively.

---

## Conventions

### Response envelope

Every response uses the same envelope, so a client can parse success and failure the same way:

```json
{
  "status": "success",
  "message": "Students created successfully",
  "data": {},
  "meta": {},
  "details": {}
}
```

`message` is omitted on plain `200 OK` responses and present on `201 Created` and errors. `meta` appears on paginated responses, `details` on validation failures.

### Pagination

List endpoints accept `page` and `limit` query parameters and report back the values that were actually applied — each resource enforces its own sensible default and an upper bound (100 for standard collections), so a client never has to guess what it received:

```json
{
  "status": "success",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "count": 10,
    "total_items": 42,
    "total_pages": 5
  }
}
```

### Errors

```json
{
  "status": "error",
  "message": "Validation failed",
  "details": {
    "Email": "email",
    "FirstName": "required"
  }
}
```

| Status | Used for                                               |
| ------ | ------------------------------------------------------ |
| `400`  | Malformed JSON or invalid input                        |
| `401`  | Missing, malformed or expired token; wrong credentials |
| `403`  | Authenticated but not permitted; inactive account      |
| `404`  | Resource does not exist                                |
| `422`  | Validation failed — field-level details included       |
| `429`  | Rate limit exceeded                                    |
| `500`  | Unexpected server error                                |

### Bulk operations

Bulk student endpoints take an array and answer with a batch summary:

```json
{
  "status": "success",
  "message": "Students created successfully",
  "data": {
    "items": [],
    "meta": { "count": 2, "total": 2 }
  }
}
```

Limits are enforced by validation (`2,000` per create, `1,000` per update, `100` per delete), and writes run inside a single transaction so a batch is all-or-nothing.

---

## Security

- **Password hashing** — Argon2id with a 16-byte random salt and 32-byte key (64 MB memory, 1 iteration, 4 threads), stored in standard PHC format (`$argon2id$v=19$m=65536,t=1,p=4$...`). Verification uses `subtle.ConstantTimeCompare` to avoid timing leaks.
- **Never in the response** — password hashes, reset tokens and token expiry are `json:"-"`, and users are serialised through a dedicated `UserResponse` type so a hash cannot reach a client by accident.
- **JWT handling** — the signing method is pinned to HMAC on parse, which blocks algorithm-substitution attacks, and the issuer, subject, `iat` and `exp` claims are all set explicitly.
- **One-time tokens** — activation and reset tokens are 32 bytes from `crypto/rand`, hex encoded, stored with an expiry, and cleared the moment they are used.
- **Rate limiting** — login, forgot-password, reset-password and activate each get their own per-IP token bucket, and idle entries are garbage-collected after 30 minutes.
- **No user enumeration** — forgot-password answers identically whether or not the address is registered, and login failures never distinguish an unknown e-mail from a wrong password.
- **Role gates** — administrative user management sits behind `RequireRole(super_admin)` at the router level, so a handler cannot be reached without the role by accident.
- **No plaintext in transit or at rest** — passwords are hashed before they touch the repository, and secrets come from the environment, never from source.

---

## Testing

```bash
go test ./...
```

The test suite covers the rate limiter, including the case that matters most — that repeated requests from the same IP reuse a single bucket rather than minting a fresh one and quietly defeating the limit. The rest of the codebase is structured for testability (interfaces in the domain layer, dependency injection through constructors), so services and handlers can be exercised against fakes without a database.

---

## License

Released under the MIT License — see [LICENSE](https://github.com/ednanf/school-api/blob/master/LICENSE).

Copyright © 2026 Ednan R. Frizzera Filho.
