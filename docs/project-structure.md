# Project Structure

This document describes the folder layout of the Go Hiring Challenge project and the purpose of each directory.

---

## Root files

| File | Purpose |
|---|---|
| `.env` | Environment variables (DB credentials, ports, etc.) |
| `README.md` | Project overview and setup instructions |
| `Makefile` | Useful commands: `run`, `seed`, `test`, `docker-up/down`, `tidy` |
| `docker-compose.yml` | Docker services definition (PostgreSQL, env vars) |

---

## Folders

### `cmd/`
Entry points of the project.

- **`cmd/server/`** — Wires up the database, repositories, and HTTP handlers, then starts the REST API server with graceful shutdown support.
- **`cmd/seed/`** — Reads all `.sql` files with migrations and seeds to execute them in ethe db.

---

### `app/`
Application layer: HTTP handlers, DTOs, and infrastructure helpers.

- **`app/api/`** — Shared HTTP response helpers used by all handlers to write consistent JSON responses.
- **`app/catalog/`** — Handler for the `/catalog` endpoints. Supports listing products with offset pagination and filtering by category or price, plus a product detail endpoint at `/catalog/{code}`.
- **`app/categories/`** — Handler for the `/categories` endpoints. Supports listing all categories and creating new ones.
- **`app/database/`** — PostgreSQL connection factory using GORM. Returns a db instance and a close function.
- **`app/dto/`** — Data Transfer Objects that shape the JSON request/response payloads, keeping the HTTP layer decoupled from the internal models.

---

### `models/`
Domain models and repository implementations.

- **`models/`** (root files) — GORM model structs and their repository interfaces and implementations that interact with the database.
- **`models/mocks/`** — Mock implementations of the repository interfaces, used in unit tests to avoid real database calls.

---

### `sql/`
SQL migration and seed scripts. The seed command runs these files in filename order to create tables and populate initial data.
