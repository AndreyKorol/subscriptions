# Subscriptions Service

REST API сервис для управления подписками и агрегации их стоимости. Написан на Go, использует PostgreSQL.

## Технологии

- **Go 1.26** — стандартная библиотека `net/http` (маршрутизация с методами), `log/slog` (структурированное логирование)
- **PostgreSQL** — драйвер [`pgx/v5`](https://github.com/jackc/pgx)
- **Валидация** — [`go-playground/validator/v10`](https://github.com/go-playground/validator)
- **SQL Builder** — [`huandu/go-sqlbuilder`](https://github.com/huandu/go-sqlbuilder)
- **Декодирование query params** — [`gorilla/schema`](https://github.com/gorilla/schema)
- **Типизированные ошибки API** — собственный пакет `internal/errs` (JSON-формат `{ "error": { "code", "message", "details" } }`)
- **Контейнеризация** — Docker + Docker Compose
- **OpenAPI 3.1** — спецификация в `openapi/openapi.yaml`

## Быстрый старт

```bash
docker compose up --build
```

- Сервис: `http://localhost:8080`
- Swagger UI (интерактивная документация): `http://localhost:8081`

Миграции применяются автоматически через контейнер `migrate`.

## Конфигурация

| Переменная         | По умолчанию    | Описание              |
|--------------------|-----------------|-----------------------|
| `POSTGRES_USER`    | `postgres`      | Пользователь БД       |
| `POSTGRES_PASSWORD`| `postgres`      | Пароль БД             |
| `POSTGRES_DB`      | `subscriptions` | Название БД           |
| `DB_HOST`          | `postgres`      | Хост БД               |
| `DB_PORT`          | `5432`          | Порт БД               |
| `DB_SSLMODE`       | `disable`       | Режим SSL             |
| `APP_PORT`         | `8080`          | Порт HTTP-сервера     |
| `APP_ENV`          | `development`   | Окружение             |
| `APP_LOG_LEVEL`    | `debug`         | Уровень логирования   |
| `SWAGGER_PORT`     | `8081`          | Порт Swagger UI       |

## API Endpoints

| Метод    | Путь                   | Описание                          |
|----------|------------------------|-----------------------------------|
| `GET`    | `/subscriptions`       | Список подписок (с фильтрацией)   |
| `POST`   | `/subscriptions`       | Создать подписку                  |
| `GET`    | `/subscriptions/{id}`  | Получить подписку по ID           |
| `PATCH`  | `/subscriptions/{id}`  | Обновить подписку                 |
| `DELETE` | `/subscriptions/{id}`  | Удалить подписку                  |
| `GET`    | `/subscriptions/agg`   | Агрегировать стоимость подписок   |

### Query-параметры (для Index и Aggregate)

| Параметр      | Тип     | Формат       | Описание                  |
|---------------|---------|--------------|---------------------------|
| `service_name`| string  |              | Фильтр по названию сервиса|
| `user_id`     | string  | UUID v4      | Фильтр по пользователю    |
| `start_date`  | string  | `MM-YYYY`    | Начало периода (включит.) |
| `end_date`    | string  | `MM-YYYY`    | Конец периода (включит.)  |

### Формат ответов

- Одиночный объект: `{ "data": { ... } }`
- Коллекция: `{ "data": { "items": [ ... ] } }`
- Агрегация: `{ "data": { "sum_price": 123 } }`
- Ошибка: `{ "error": { "code": "BAD_REQUEST", "message": "..." } }`
- Ошибка с деталями: `{ "error": { "code": "BAD_REQUEST", "message": "...", "details": [{"field": "price", "message": "..."}] } }`

### Коды ошибок

| HTTP | Код             | Описание                          |
|------|-----------------|-----------------------------------|
| 400  | `BAD_REQUEST`   | Валидация или неверный запрос     |
| 404  | `NOT_FOUND`     | Ресурс не найден                  |
| 500  | `INTERNAL_ERROR`| Внутренняя ошибка (без деталей)   |

### Примеры

**Создать подписку**

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H 'Content-Type: application/json' \
  -d '{
    "service_name": "Spotify",
    "price": 199,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "01-2025"
  }'
```

**Список подписок с фильтрацией**

```bash
curl 'http://localhost:8080/subscriptions?user_id=550e8400-e29b-41d4-a716-446655440000&start_date=01-2025&end_date=12-2025'
```

**Агрегация стоимости**

```bash
curl 'http://localhost:8080/subscriptions/agg?user_id=550e8400-e29b-41d4-a716-446655440000'
```

**Получить подписку по ID**

```bash
curl http://localhost:8080/subscriptions/1
```

**Обновить подписку**

```bash
curl -X PATCH http://localhost:8080/subscriptions/1 \
  -H 'Content-Type: application/json' \
  -d '{"price": 299}'
```

**Удалить подписку**

```bash
curl -X DELETE http://localhost:8080/subscriptions/1
```

## Архитектура

```
HTTP Request
  ↓
LoggingMiddleware (логирование метода, пути, статуса, времени)
  ↓
Controller (парсинг, валидация, JSON)
  ↓
Service (бизнес-логика)
  ↓
Repository (SQL-запросы через pgx + go-sqlbuilder)
  ↓
PostgreSQL
```

## Структура проекта

```
├── subscriptions.go             # Точка входа (main)
├── internal/
│   ├── config/config.go         # Загрузка конфигурации из env
│   ├── errs/
│   │   ├── errors.go            # Тип Error и конструкторы (NotFound, BadRequest, Internal)
│   │   └── validation.go        # Трансляция validator.ValidationErrors в JSON
│   ├── controllers/
│   │   ├── subscriptions.go     # HTTP-обработчики
│   │   ├── middleware.go        # Логирующая middleware
│   │   └── requests.go          # DTO для запросов
│   ├── migrations/              # SQL-миграции
│   ├── models/subscription.go   # Модели данных
│   ├── repos/subscription.go    # Слой доступа к данным
│   └── services/
│       ├── subscription.go      # Бизнес-логика
│       └── manager.go           # DI-контейнер сервисов
├── openapi/openapi.yaml         # OpenAPI 3.1 спецификация
├── docker-compose.yaml          # Docker Compose (app + postgres + migrate + swagger-ui)
├── Dockerfile                   # Сборка приложения
└── .env                         # Локальные переменные окружения
```

## Разработка (локально без Docker)

```bash
# Убедитесь, что PostgreSQL запущен и создана БД
go run subscriptions.go
```

## OpenAPI / Swagger документация

Файл `openapi/openapi.yaml` содержит полное описание API в формате OpenAPI 3.1.

При запуске через Docker Compose интерактивная документация доступна через Swagger UI:
`http://localhost:8081`

Для генерации статической HTML-документации (Redoc):
```bash
npx @redocly/cli build-docs openapi/openapi.yaml -o openapi/openapi.html
```
