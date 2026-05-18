# Messenger

Простой мессенджер на Go + PostgreSQL с REST API, WebSocket и Vue/Vuetify web-клиентом.

Backend сам инициализирует базу данных при старте, если включён `INIT_DB=true`.

## Стек

- Go
- PostgreSQL
- REST API
- WebSocket
- Vue 3 + Vuetify
- Docker Compose для локальной PostgreSQL

## Запуск backend

Сначала поднять PostgreSQL:

```bash
docker compose up -d
```

Затем запустить сервер:

```bash
go run ./cmd/server
```

По умолчанию backend слушает:

```text
http://localhost:8080
```

## Запуск webclient

```bash
cd webclient
npm install
npm run dev
```

## Переменные окружения backend

```text
DATABASE_URL
ADDR
INIT_DB
```

### DATABASE_URL

Строка подключения к PostgreSQL.

Значение по умолчанию:

```text
postgres://postgres:postgres@localhost:5432/messenger?sslmode=disable
```

### ADDR

Адрес HTTP-сервера.

Значение по умолчанию:

```text
:8080
```

### INIT_DB

Запускать SQL-инициализацию базы данных при старте.

Значение по умолчанию:

```text
true
```

SQL-схема находится здесь:

```text
internal/repository/postgres/schema.sql
```

## Структура проекта

```text
cmd/server/main.go                  # точка входа backend, env, wiring зависимостей
internal/handler/                   # HTTP handlers и WebSocket handler
internal/service/                   # бизнес-логика и валидация
internal/repository/                # интерфейсы репозиториев
internal/repository/postgres/       # PostgreSQL-реализация, SQL-запросы и schema.sql
internal/model/                     # модели JSON/DB
internal/middleware/                # logging, auth, CORS middleware
webclient/                          # Vue/Vuetify frontend
index.html                          # простой standalone HTML-клиент
docker-compose.yml                  # PostgreSQL для локального запуска
```

## База данных

При старте backend создаёт таблицы, если они ещё не существуют:

- `users`
- `channels`
- `posts`
- `memberships`
- `tokens`

Также создаётся расширение PostgreSQL:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

UUID генерируются на стороне PostgreSQL через `gen_random_uuid()`.

## Авторизация

В проекте есть token-based авторизация.

Текущий middleware авторизации пропускает без проверки:

- `POST /auth`
- WebSocket URL, начинающиеся с `/ws?`

Остальные REST-запросы требуют `Authorization: Bearer <token>`, включая `GET /health` и `POST /users`.

Это значит, что для пустой базы нужен bootstrap первого пользователя: например, временно отключить auth middleware на время разработки, создать пользователя напрямую в базе или добавить `POST /users` в список публичных endpoint'ов.

После появления пользователя token можно получить через:

```http
POST /auth
```

Защищённые REST-запросы должны передавать token в заголовке:

```http
Authorization: Bearer <token>
```

WebSocket передаёт token через query-параметр:

```text
ws://localhost:8080/ws?channel_id=<channel_id>&token=<token>
```

## Формат ошибок

Большинство ошибок возвращаются в JSON:

```json
{
  "error": "error message"
}
```

Если запрос отклонён middleware авторизации, ответ может быть обычным текстом:

```text
invalid authorization header
```

## REST API

Базовый URL:

```text
http://localhost:8080
```

### Healthcheck

```http
GET /health
```

Возвращает:

```json
{
  "status": "ok"
}
```

Важно: текущий auth middleware защищает `/health`, поэтому для запроса нужен заголовок `Authorization`.

### Создать пользователя

```http
POST /users
```

Headers:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "name": "Slava",
  "email": "slava@example.com",
  "password_hash": "password-or-hash"
}
```

Response:

```json
{
  "id": "uuid",
  "name": "Slava",
  "email": "slava@example.com",
  "password_hash": "password-or-hash"
}
```

### Получить пользователей

```http
GET /users
```

Response:

```json
[
  {
    "id": "uuid",
    "name": "Slava",
    "email": "slava@example.com",
    "password_hash": ""
  }
]
```

### Авторизация

```http
POST /auth
```

Request:

```json
{
  "email": "slava@example.com",
  "pwd": "password-or-hash"
}
```

Response:

```json
"generated-token"
```

Этот token нужно использовать в `Authorization: Bearer <token>`.

### Создать канал

```http
POST /channels
```

Request:

```json
{
  "type": "public"
}
```

Допустимые значения `type`:

```text
direct
public
private
```

Response:

```json
{
  "id": "uuid",
  "type": "public"
}
```

### Получить каналы

```http
GET /channels
```

Response:

```json
[
  {
    "id": "uuid",
    "type": "public"
  }
]
```

### Добавить пользователя в канал

```http
POST /channels/{channel_id}/members
```

Request:

```json
{
  "user_id": "uuid"
}
```

Response:

```json
{
  "user_id": "uuid",
  "channel_id": "uuid"
}
```

### Получить участников канала

```http
GET /channels/{channel_id}/members
```

Response:

```json
[
  {
    "id": "uuid",
    "name": "Slava",
    "email": "slava@example.com",
    "password_hash": ""
  }
]
```

### Создать сообщение

```http
POST /channels/{channel_id}/posts
```

Headers:

```http
Authorization: Bearer <token>
```

Request:

```json
{
  "text": "hello"
}
```

Response:

```json
{
  "id": "uuid",
  "text": "hello",
  "channel_id": "uuid",
  "author": "uuid",
  "created_at": "2026-05-18T12:00:00Z"
}
```

После создания сообщение рассылается всем WebSocket-клиентам, подключённым к этому каналу.

### Получить сообщения канала

```http
GET /channels/{channel_id}/posts
```

Response:

```json
[
  {
    "id": "uuid",
    "text": "hello",
    "channel_id": "uuid",
    "author": "uuid",
    "created_at": "2026-05-18T12:00:00Z"
  }
]
```

## WebSocket API

Endpoint:

```text
ws://localhost:8080/ws?channel_id=<channel_id>&token=<token>
```

Query-параметры:

- `channel_id` — UUID канала
- `token` — token, полученный через `POST /auth`

Отправка сообщения:

```json
{
  "text": "hello from websocket"
}
```

Получение сообщения:

```json
{
  "id": "uuid",
  "text": "hello from websocket",
  "channel_id": "uuid",
  "author": "uuid",
  "created_at": "2026-05-18T12:00:00Z"
}
```

Если сообщение невалидное, сервер отправляет:

```json
{
  "error": "text is required"
}
```

## CORS

Backend содержит CORS middleware.

Разрешены методы:

```text
GET, POST, OPTIONS
```

Разрешены заголовки:

```text
Content-Type, Authorization
```

Это нужно для работы frontend-клиентов, которые обращаются к backend с другого origin.

## Примеры curl

### Создать пользователя

В текущей конфигурации этот запрос требует уже существующий token.

```bash
curl -X POST http://localhost:8080/users \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Slava","email":"slava@example.com","password_hash":"123"}'
```

### Получить token

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth \
  -H 'Content-Type: application/json' \
  -d '{"email":"slava@example.com","pwd":"123"}')
```

### Создать канал

```bash
curl -X POST http://localhost:8080/channels \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"type":"public"}'
```

### Получить каналы

```bash
curl http://localhost:8080/channels \
  -H "Authorization: Bearer $TOKEN"
```

### Добавить пользователя в канал

```bash
curl -X POST http://localhost:8080/channels/<channel_id>/members \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"<user_id>"}'
```

### Создать сообщение

```bash
curl -X POST http://localhost:8080/channels/<channel_id>/posts \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"text":"hello"}'
```

### Получить сообщения

```bash
curl http://localhost:8080/channels/<channel_id>/posts \
  -H "Authorization: Bearer $TOKEN"
```
