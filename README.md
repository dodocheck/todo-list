# Todo List (Check-List) — pet-проект на Go

Небольшая **микросервисная система для управления задачами**: внешний HTTP API проксирует запросы в gRPC DB-сервис, который работает с PostgreSQL и кэширует данные в Redis. Действия пользователя дополнительно публикуются в Kafka и пишутся в файл отдельным logger-сервисом.

![Схема проекта](assets/image.png)

> Исходное ТЗ и чек-лист выполнения — в [TODO.md](TODO.md).

## Возможности

- CRUD для задач: **создать / получить список / отметить выполненной / удалить**
- Микросервисы:
  - **api-service** — HTTP API (Gorilla/mux) + продюсер событий в Kafka
  - **db-service** — gRPC API + PostgreSQL, Redis-кэш с TTL и инвалидацией
  - **logger-service** — Kafka consumer, пишет события в лог
- Всё поднимается через **Docker Compose** (Postgres, Redis, Kafka + сервисы)
- **Unit** и **интеграционные** тесты (happy path end-to-end)

> Нюанс: для простоты демо `db-service` при старте пересоздаёт таблицу `tasks` и заполняет её тестовыми задачами.

## Быстрый старт

Требования: Docker + Docker Compose.

```bash
make deploy
# API будет доступно на http://localhost:9089 (значения по умолчанию — в deployment/.env)
```

Остановить и удалить volumes:

```bash
make down
```

## HTTP API

### `POST /create` — создать задачу

**Body:**

```json
{"title":"...","text":"..."}
```

**Ответ:** `201 Created` → созданная задача

---

### `GET /list` — получить список задач

**Ответ:** `200 OK` → список задач

---

### `PUT /done` — отметить задачу выполненной

**Body:**

```json
{"Id":1}
```

**Ответ:** `200 OK` → обновлённая задача

---

### `DELETE /delete` — удалить задачу

**Body:**

```json
{"Id":1}
```

**Ответ:** `204 No Content`

---

### Примеры запросов

```bash
curl -X POST http://localhost:9089/create \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk","text":"2 liters"}'

curl http://localhost:9089/list

curl -X PUT http://localhost:9089/done \
  -H 'Content-Type: application/json' \
  -d '{"Id":1}'

curl -X DELETE http://localhost:9089/delete \
  -H 'Content-Type: application/json' \
  -d '{"Id":1}'
```

## Разработка

* `make test` — unit-тесты
* `make integration-test` — интеграционные тесты (Docker Compose + `-tags=integration`)
* `make lint` — golangci-lint
* `make proto-gen` — генерация gRPC/Protobuf в `pkg/pb` (нужны `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`)

```

Если хочешь — могу также привести `TODO.md` к такому же “github-friendly” стилю (чтобы там тоже не ломались списки/чек-боксы/заголовки).
```
