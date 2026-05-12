# Messenger backend

Минимальный backend для мессенджера на Go + PostgreSQL: REST API, WebSocket и автоматический init базы при старте.

## Запуск

```bash
docker compose up -d
go run ./cmd/server
```

По умолчанию приложение подключается к:

```text
postgres://postgres:postgres@localhost:5432/messenger?sslmode=disable
```

Переменные окружения:

- `DATABASE_URL` - строка подключения к PostgreSQL
- `ADDR` - адрес HTTP-сервера, по умолчанию `:8080`
- `INIT_DB` - запускать SQL-схему из `internal/repository/postgres/schema.sql` при старте, по умолчанию `true`

## Структура

```text
cmd/server/main.go                  # запуск, env, wiring зависимостей
internal/handler/                   # REST + WebSocket
internal/service/                   # простая бизнес-логика и валидация
internal/repository/                # интерфейсы репозиториев
internal/repository/postgres/       # SQL и PostgreSQL-реализация
internal/model/                     # модели JSON/DB
internal/middleware/                # middleware
```

## REST API

```http
GET /health
GET /users
POST /users
GET /channels
POST /channels
GET /channels/{channel_id}/members
POST /channels/{channel_id}/members
GET /channels/{channel_id}/posts
POST /channels/{channel_id}/posts
```

Примеры:

```bash
curl -X POST http://localhost:8080/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"Slava","email":"slava@example.com"}'

curl -X POST http://localhost:8080/channels \
  -H 'Content-Type: application/json' \
  -d '{"type":"public"}'

curl -X POST http://localhost:8080/channels/<channel_id>/members \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"<user_id>"}'

curl -X POST http://localhost:8080/channels/<channel_id>/posts \
  -H 'Content-Type: application/json' \
  -d '{"text":"hello"}'
```

## WebSocket

Подключение:

```text
ws://localhost:8080/ws?channel_id=<channel_id>
```

Отправка сообщения:

```json
{"text":"hello from ws"}
```

Сервер сохраняет сообщение в `posts` и рассылает новый post всем WebSocket-клиентам этого канала.
