# Penaut Messenger

Penaut Messenger is a small messenger project built with Go, PostgreSQL, MinIO, WebSocket, and a Vue/Vuetify frontend.

> Status: this project is still in development. APIs, auth behavior, database schema, and frontend flows may change.

## Stack

- Go backend
- PostgreSQL
- MinIO object storage
- REST API
- WebSocket messaging
- Vue 3 + Vuetify frontend
- Docker Compose for local development

## Project Structure

```text
cmd/server/                    Backend entrypoint
internal/handler/              HTTP and WebSocket handlers
internal/service/              Business logic
internal/repository/           Repository interfaces
internal/repository/postgres/  PostgreSQL implementation and schema
internal/storage/minio/        MinIO storage integration
internal/middleware/           Auth, CORS, and logging middleware
webclient/                     Vue/Vuetify frontend
docker-compose.yml             Local development stack
```

## Run With Docker Compose

```bash
docker compose up --build
```

Default public services:

```text
Backend:  http://localhost:8080
Frontend: http://localhost:3000
```

PostgreSQL and MinIO are used internally by the compose network. Their public port mappings are currently disabled in `docker-compose.yml`.

## Run Backend Locally

Start dependencies first:

```bash
docker compose up postgres minio
```

Then run the backend:

```bash
go run ./cmd/server
```

The backend initializes the database schema on startup when `INIT_DB=true`.

## Run Frontend Locally

```bash
cd webclient
npm install
npm run dev
```

The frontend dev server listens on:

```text
http://localhost:3000
```

## Environment Variables

All variables have defaults in `docker-compose.yml`, so the project can start without a `.env` file.

Common variables:

```text
DOCKER_BE_PORT=8080
DOCKER_FE_PORT=3000

DOCKER_DB_NAME=postgres
DOCKER_DB_PORT=5432
DB_NAME=messenger
DB_USER=postgres
DB_PASSWORD=postgres

DOCKER_MINIO_NAME=minio
DOCKER_MINIO_PORT_WEB=9000
DOCKER_MINIO_PORT_API=9001
MINIO_ACCESS_KEY=minio
MINIO_SECRET_KEY=minio123
MINIO_BUCKET=files

INIT_DB=true
```

## API Overview

Base URL:

```text
http://localhost:8080
```

Main REST endpoints:

```http
GET  /health
POST /auth

GET  /users
POST /users
GET  /users/{user_id}/pic
POST /users/{user_id}/direct

GET  /channels
POST /channels
GET  /channels/{channel_id}/members
POST /channels/{channel_id}/members
GET  /channels/{channel_id}/posts
POST /channels/{channel_id}/posts
GET  /channels/{channel_id}/pic

GET  /posts/{post_id}/attachments
GET  /attachment/{file_path}
```

WebSocket endpoint:

```text
ws://localhost:8080/ws?channel_id=<channel_id>&token=<token>
```

## Authentication

The backend uses token-based authentication.

Login:

```http
POST /auth
```

Request:

```json
{
  "email": "user@example.com",
  "pwd": "password-or-hash"
}
```

Response:

```json
"generated-token"
```

Protected REST requests use:

```http
Authorization: Bearer <token>
```

WebSocket requests pass the token as a query parameter:

```text
ws://localhost:8080/ws?channel_id=<channel_id>&token=<token>
```

Development note: auth and user bootstrap are still being refined. In a fresh database, creating the first user may require a temporary development bootstrap step.

## Database

The PostgreSQL schema is defined in:

```text
internal/repository/postgres/schema.sql
```

The current schema includes:

- users
- channels
- posts
- memberships
- tokens
- attachments/storage-related tables as the project evolves

## Development Checks

Run Go package checks:

```bash
go test ./...
```

Validate compose configuration:

```bash
docker compose config
```

## Notes

- The project is not production-ready yet.
- Error formats and auth rules may still change.
- Docker Compose waits for PostgreSQL health before starting the backend.
- CORS middleware is enabled for frontend/backend development.
