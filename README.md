# **Go High-Load Microservice**

### CRUD API • Prometheus • Rate Limiting • MinIO • Docker

Проект представляет собой учебный высоконагруженный микросервис на языке **Go**, разработанный в рамках практического задания.
Сервис реализует CRUD API, конкурентную обработку операций, метрики Prometheus, rate limiting и интеграцию с S3-совместимым хранилищем MinIO.


# **Возможности**

### CRUD-операции над пользователями

* `GET /api/users`
* `GET /api/users/{id}`
* `POST /api/users`
* `PUT /api/users/{id}`
* `DELETE /api/users/{id}`

### Асинхронная обработка в goroutines

* audit log (канал + отдельная горутина)
* уведомления (stub sender)

### Rate Limiter

* GET — без ограничения
* POST/PUT/DELETE — лимит 1000 RPS

### Prometheus-метрики

* RPS по методам и эндпоинтам
* Latency (Histogram)
* Endpoint `/metrics`

### Интеграция с MinIO

* Загрузка тестового файла в S3-бакет
* Эндпоинт `/api/integration/upload-test`

### Полная контейнеризация

* Dockerfile (multi-stage build)
* docker-compose (Go + MinIO)


# **Быстрый старт**

## 1. Клонирование проекта

```bash
git clone https://github.com/your-repo/go-microservice.git
cd go-microservice
```

## 2. Запуск через Docker Compose

```bash
docker compose up --build
```

После запуска доступны:

| Компонент          | URL                                                                |
| ------------------ | ------------------------------------------------------------------ |
| API                | [http://localhost:8080/api/users](http://localhost:8080/api/users) |
| Prometheus metrics | [http://localhost:8080/metrics](http://localhost:8080/metrics)     |
| MinIO console      | [http://localhost:9001](http://localhost:9001)                     |
| MinIO S3 endpoint  | [http://localhost:9000](http://localhost:9000)                     |


# **Структура проекта**

```
go-microservice/
├── main.go                  # Точка входа
├── handlers/                # HTTP-обработчики
├── services/                # Бизнес-логика
├── storage/                 # Потокобезопасное хранилище пользователей
├── utils/                   # Логирование, rate limiting, валидация
├── metrics/                 # Prometheus метрики
├── Dockerfile
├── docker-compose.yml
└── README.md
```


# **Rate Limiting**

Используется пакет `golang.org/x/time/rate`.

### Принцип:

* для операций чтения (`GET`) ограничение **не применяется**, чтобы обеспечить максимальный RPS;
* для операций изменений (`POST`, `PUT`, `DELETE`) используется лимит:

```go
limiter := rate.NewLimiter(1000, 5000)
```

Это гарантирует защиту сервиса, сохраняя высокую пропускную способность чтения.


# **Метрики Prometheus**

Доступны по адресу:

```
GET /metrics
```

Метрики включают:

* `http_requests_total` — общее количество HTTP-запросов
* `http_request_duration_seconds` — гистограмма задержек

Интегрируются в Prometheus/Grafana без дополнительной настройки.


# **Интеграция с MinIO**

MinIO запускается в составе `docker-compose` и доступен по адресу:

* Console: [http://localhost:9001](http://localhost:9001)
* S3: [http://localhost:9000](http://localhost:9000)

Пример эндпоинта, создающего тестовый файл в бакете:

```
POST /api/integration/upload-test
```


# **Нагрузочное тестирование**

Использовался инструмент `wrk`.

Стресс-тест:

```
wrk -t12 -c500 -d60s http://localhost:8080/api/users
```

Производительный тест:

```
wrk -t12 -c200 -d60s http://localhost:8080/api/users
```

### Итоги:

* RPS: **≈38–40k**
* Ошибки: **0%**
* Latency: **≈5.7 ms (<10 ms)**
* Сервис стабильно работает в Docker даже в среде Windows → VMware → Ubuntu


# **Переменные окружения**

Можно задать в `.env` или в `docker-compose.yml`:

| Переменная       | Описание | Значение по умолчанию |
| ---------------- | -------- | --------------------- |
| MINIO_ENDPOINT   | S3 host  | minio:9000            |
| MINIO_ACCESS_KEY | логин    | minioadmin            |
| MINIO_SECRET_KEY | пароль   | minioadmin            |
| MINIO_BUCKET     | S3-бакет | go-microservice       |


# **Запуск без Docker (локально)**

```bash
go mod download
go run main.go
```


# **Тестирование API**

Получить всех пользователей:

```bash
curl http://localhost:8080/api/users
```

Создать пользователя:

```bash
curl -X POST http://localhost:8080/api/users \
  -d '{"name":"Alice","email":"alice@example.com"}' \
  -H "Content-Type: application/json"
```


# **Требования окружения**

* Go 1.22+
* Docker 24+
* docker-compose 2+


# **Лицензия**

MIT