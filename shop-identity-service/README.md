# Identity Service --- Golang Microservice

A **production-style Identity & Authentication microservice** built with
Go, designed to act as the **central authentication authority** in a
microservices architecture.

This service provides:

-   Secure **user registration & login**
-   **JWT-based authentication**
-   **gRPC internal identity API** for other microservices
-   **PostgreSQL persistence**
-   **Dockerized local development environment**
-   Clean, extensible **Hexagonal / Clean Architecture**

This project is built as part of a **microservices platform** intended
to include:

-   API Gateway
-   Product Service
-   Inventory Service
-   Order & Payment Services (future)

------------------------------------------------------------------------

## ✅ Architecture Overview

This service follows **Clean Architecture / Hexagonal Architecture**
principles:

    cmd/
      └── identity-service/        → Application entrypoint

    internal/
      ├── config/                 → Environment-based configuration
      ├── httpserver/             → HTTP server bootstrap
      └── identity/               → Identity business module
           ├── domain/            → Core business entities & validation
           ├── repository/        → Persistence abstraction (interface + GORM implementation)
           ├── transport/
           │    ├── http/         → REST API (public)
           │    └── grpc/         → gRPC API (internal)
           ├── jwt.go             → JWT generation & verification
           ├── password.go        → Password hashing (bcrypt)
           ├── service.go         → Business logic
           ├── context.go         → Claims injection into context
           └── types.go           → Public type re-exports

    infrastructure/
      ├── Dockerfile
      └── docker-compose.yml

### Key Design Principles

-   **Domain is pure** (no HTTP, no DB, no GORM)
-   **Service layer orchestrates business logic**
-   **Repository is an interface** (infrastructure implements it)
-   **JWT & gRPC are treated as infrastructure adapters**
-   **HTTP & gRPC are independent transports**
-   **Everything is injectable & testable**

------------------------------------------------------------------------

## ✅ Features

### Authentication

-   User Registration
-   User Login
-   Password hashing with **bcrypt**
-   JWT generation (HS256)

### Authorization

-   JWT verification middleware
-   Claims injection into request context
-   Protected HTTP endpoint (`/me`)

### Persistence

-   PostgreSQL via **GORM**
-   Domain ↔ DB model separation
-   Repository abstraction (`interface`)

### Internal Communication

-   gRPC Identity API:
    -   `GetUser`
    -   `ValidateToken`
-   Used by other microservices (Product, Inventory, Gateway)

### Infrastructure

-   Dockerized service
-   Docker Compose for local Postgres
-   Environment-based configuration

------------------------------------------------------------------------

## ✅ Tech Stack

-   Go 1.22+
-   Chi Router (HTTP)
-   gRPC (internal APIs)
-   GORM + PostgreSQL
-   JWT (golang-jwt v5)
-   bcrypt
-   Docker & Docker Compose

------------------------------------------------------------------------

## ✅ REST API

### Register User

``` http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "vinh",
  "password": "secure123"
}
```

### Login

``` http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure123"
}
```

### Get Current User (JWT Protected)

``` http
GET /api/v1/auth/me
Authorization: Bearer <access_token>
```

------------------------------------------------------------------------

## ✅ gRPC API (Internal)

``` proto
rpc GetUser(GetUserRequest) returns (GetUserResponse);
rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
```

Used internally by: - API Gateway - Product Service - Inventory Service

------------------------------------------------------------------------

## ✅ Event-Driven Welcome Emails

-   After every successful registration, the identity service publishes a
    `user_created` event to Kafka via the REST Proxy.
-   A new `shop-email-service` (see `../shop-email-service`) consumes the
    topic and sends a welcome email (currently logged to stdout via a
    pluggable `Mailer` interface).
-   The event payload is defined in `../platform-events/pkg/events/` so additional services
    can reuse the same contract.

Run the email service locally:

``` bash
cd ../shop-email-service
cp .env.example .env
make run
```

------------------------------------------------------------------------

## ✅ Environment Configuration

Create `.env` from `.env.example`:

``` env
HTTP_PORT=8081
GRPC_PORT=9091

JWT_SECRET=dev-secret-changeme
JWT_ISSUER=identity-service
JWT_EXPIRES_IN=15m

DB_DSN=postgres://postgres:postgres@localhost:5432/identity?sslmode=disable
KAFKA_REST_URL=http://localhost:8082
KAFKA_TOPIC_USER_CREATED=user_created
```

------------------------------------------------------------------------

## ✅ Run Locally (Docker)

``` bash
make dev
```

or

``` bash
docker compose up --build
```

### Services Available Locally

| Service             | Address               |
| ------------------- | --------------------- |
| Identity HTTP       | http://localhost:8081 |
| Identity gRPC       | localhost:9091        |
| PostgreSQL          | localhost:5432        |
| Kafka REST Proxy    | http://localhost:8082 |
| Kafka Broker        | localhost:19092       |

------------------------------------------------------------------------

## ✅ Run Without Docker

``` bash
go run ./cmd/identity-service
```

(Requires local PostgreSQL)

------------------------------------------------------------------------

## ✅ Testing

``` bash
make test
```

Planned: - Service unit tests - JWT unit tests - HTTP handler tests -
PostgreSQL integration tests (Testcontainers)

------------------------------------------------------------------------

## ✅ Security Model

-   **User authentication** → JWT (HS256)
-   **Public APIs** → secured by JWT middleware
-   **Internal gRPC** → trusted network model (K8s / private VPC)
-   Future upgrades:
    -   mTLS for internal gRPC
    -   RS256 JWT signing
    -   Role-based access control (RBAC)

------------------------------------------------------------------------


## ✅ Project Goals

This project is a **hands-on learning and skill-deepening side project**, built to:

- Practice building a **real-world microservice architecture** from scratch
- Apply **Clean Architecture principles in Go** in a production-style codebase
- Gain deep practical experience with **JWT-based authentication**
- Learn and implement **service-to-service communication using gRPC**
- Design **database abstractions using the repository pattern**
- Develop and operate **containerized services using Docker & Docker Compose**
- Prepare for **Kubernetes-based deployment and cloud-native systems**

The goal is not only to build a working system, but to continuously **evolve the architecture**, apply best practices, and simulate how real backend platforms grow over time.

------------------------------------------------------------------------
