# Backend редактора диаграмм

Backend-сервис для [совместного редактора диаграмм](https://github.com/ilkol/diagram-editor). Представляет собой REST API для сохранения досок в PostgreSQL (`JSONB`) и планируется внедрение совместного редактирования диаграмм.

## Стек

- **Язык:** Go 1.23+
- **HTTP:** `github.com/go-chi/chi/v5`, `github.com/go-chi/cors`
- **БД:** PostgreSQL 16 (`github.com/jackc/pgx/v5/pgxpool`)
- **Миграции:** `github.com/pressly/goose/v3` (автоматически при старте)
- **Конфигурации:** `github.com/ilyakaznacheev/cleanenv`
- **Тесты:** `testing`, `net/http/httptest`, `github.com/stretchr/testify`
- **Инфраструктура:** Docker, Docker Compose

## Архитектура

```text
cmd/api/main.go          # Точка запуска приложения, подготовка зависимостей, запуск HTTP-сервера
internal/
├── config/              # Загрузчик переменных окружения и DSN-builder
├── handler/             # HTTP-обработчики
├── service/             # Бизнес-логика, валидация входных данных
├── storage/             # реализация репозиториев (PostgreSQL, pgxpool)
└── model/               # Модели, DTO и внутренние ошибки
migrations/              # SQL-миграции (goose)
```

## Схема БД

Таблица `boards`:

| Column | Type | Constraints | Description |
| :--- | :--- | :--- | :--- |
| `id` | `UUID` | `PRIMARY KEY` | Уникальный идентификатор доски |
| `title` | `TEXT` | `NOT NULL DEFAULT 'Untitled'` | Заголовок доски |
| `data` | `JSONB` | `NOT NULL DEFAULT '{}'::jsonb` | Сериализованные данные диаграммы |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | timestamp создания (UTC) |
| `updated_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT NOW()` | timestamp последнего обновления (UTC) |

## API-эндпоинты

### Проверка работы
- `GET /health` — Сервис работает (`200 OK`).

### Доски
- `POST /boards` — Создать новую доску.
  - **Тело запроса (опционально):** `{"title": "System Architecture"}`
  - **Ответ (`201 Created`):**
    ```json
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "System Architecture",
      "data": {},
      "created_at": "2026-11-25T07:00:00Z",
      "updated_at": "2026-11-25T07:00:00Z"
    }
    ```
- `GET /boards/{id}` — Запрос состояния доски по UUID (`200 OK`, `400 Bad Request`, `404 Not Found`).
- `PUT /boards/{id}` — Обновить состояние доски (`200 OK`, `400 Bad Request`, `404 Not Found`).
  - **Тело запроса:**
    ```json
    {
      "title": "Updated Architecture",
      "data": {
        "nodes": [],
        "edges": []
      }
    }
    ```

## Запуск

### Необходимые инструменты
- Go 1.23+
- Docker и Docker Compose

### Локальная сборка

1. Клонировать репозиторий:
   ```bash
   git clone https://github.com/ilkol/diagram-editor-backend.git
   cd diagram-editor-backend
   ```

2. Создать файл `.env` по примеру:
   ```bash
   cp .env.example .env
   ```

3. Запустить PostgreSQL через Docker Compose:
   ```bash
   docker compose up -d
   ```

4. Запуск приложения (миграции запустятся автоматически при подключении к БД):
   ```bash
   go run cmd/api/main.go
   ```

## Тестирование

Запуск юнит-тестов с флагом `-race`:

```bash
go test -race -v ./...
```