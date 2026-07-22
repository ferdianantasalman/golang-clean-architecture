# Golang Clean Architecture

A Go backend application following Clean Architecture principles with HTTP API (Fiber) and event-driven worker (Kafka).

## Layers

- **Entity** — Core business objects (user, contact, address, refresh token)
- **Use Case** — Business logic orchestration
- **Repository** — Database operations via GORM (MySQL)
- **Delivery** — Incoming request handlers (HTTP / Messaging)
- **Model** — Request/response DTOs and event schemas
- **Gateway** — Outbound communication to external systems (Kafka messaging)

## Tech Stack

| Category | Library |
|----------|--------|
| HTTP Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| ORM | [GORM](https://github.com/go-gorm/gorm) |
| Database | [MySQL 8](https://www.mysql.com/) |
| Cache / Auth Store | [Redis](https://redis.io/) |
| Message Broker | [Apache Kafka](https://kafka.apache.org/) (Confluent Go client) |
| Auth | JWT (access + refresh tokens) |
| Validation | [go-playground/validator](https://github.com/go-playground/validator) |
| Logger | [Logrus](https://github.com/sirupsen/logrus) |
| Config | [godotenv](https://github.com/joho/godotenv) (`.env`) |
| Test | [testify](https://github.com/stretchr/testify) |
| Migration | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Hot Reload | [air](https://github.com/air-verse/air) |

## Project Structure

```
├── api/                   # API specification
├── cmd/
│   ├── web/main.go        # HTTP server entry point
│   └── worker/main.go     # Kafka consumer entry point
├── db/migrations/         # SQL migration files
├── internal/
│   ├── config/            # App, Fiber, GORM, Redis, Kafka, Logger configs
│   ├── delivery/
│   │   ├── http/          # HTTP handlers (Fiber)
│   │   └── messaging/     # Kafka message handlers
│   ├── entity/            # Core domain entities
│   ├── gateway/messaging/ # Kafka producer
│   ├── model/             # DTOs, event schemas, converters
│   ├── repository/        # GORM repository implementations
│   ├── usecase/           # Business logic layer
│   └── util/              # Shared utilities
└── test/                  # Integration / E2E tests
```

## Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose (for MySQL, Redis, Kafka)

### Setup

```bash
# Clone and enter the project
git clone <repo-url> && cd golang-clean-architecture

# Copy environment config
cp .env.example .env

# Start infrastructure (MySQL, Redis, Kafka)
docker compose -f docker-compose.dev.yml up -d mysql redis kafka

# Run database migrations
migrate -database "mysql://root:@tcp(localhost:3306)/golang-clean-architecture?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up

# Start web server (with hot reload)
docker compose -f docker-compose.dev.yml up web

# Or run directly
go run cmd/web/main.go
```

### Run Tests

```bash
go test -v ./test/
```

## API

Full API spec available at `api/api-spec.json`.

## Database Migrations

```bash
# Create a new migration
migrate create -ext sql -dir db/migrations create_table_xxx

# Run all pending migrations
migrate -database "mysql://root:@tcp(localhost:3306)/golang-clean-architecture?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up

# Rollback
migrate -database "mysql://root:@tcp(localhost:3306)/golang-clean-architecture?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations down
```

## Architecture Flow

```
HTTP Request / Kafka Message
        │
        ▼
   Delivery (handler)
        │
        ▼
   Use Case (business logic)
        │
        ├──► Entity
        ├──► Repository ──► MySQL
        └──► Gateway ──► Kafka (external)
```

## License

MIT
