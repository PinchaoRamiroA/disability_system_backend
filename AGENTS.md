# AGENTS.md

# Project Overview

Backend for a medical leave management system built with Go using Modular Hexagonal Architecture.

Main goals:
- clean architecture
- modularity
- maintainability
- predictable structure
- AI-agent friendly development

---

# Tech Stack

## Backend
- Go 

## HTTP Framework
- Gin

## Database
- PostgreSQL 18

## ORM
- GORM

## Authentication
- JWT

## Password Hashing
- bcrypt

## Migrations
- golang-migrate

## Logging
- Zap

## Testing
- Testify

## Containerization
- Docker
- Docker Compose

---

# Architecture

This project uses:

# Modular Hexagonal Architecture (Ports & Adapters)

Flow:

```text
HTTP Request
    ↓
Handler
    ↓
UseCase
    ↓
Port Interface
    ↓
Repository
    ↓
PostgreSQL
```

Rules:
- business logic belongs in usecases
- persistence belongs in adapters/postgres
- handlers must stay thin
- modules must remain isolated

---

# Project Structure

```text
/internal
    /shared
    /modules
        /auth
        /usuarios
        /incapacidades
        /historial
        /cobros
```

Each module must contain:

```text
/domain
/usecase
/ports
/adapters
/dto
/mapper
```

---

# Layer Responsibilities

## domain
Contains:
- entities
- enums
- business rules

Must not import:
- gin
- gorm
- sql
- postgres
- http packages

---

## usecase
Contains application business logic.

Each file should represent one use case.

Examples:
- crear_usuario.go
- aprobar_incapacidad.go

---

## ports
Contains interfaces.

Use:
- input ports
- output ports

---

## adapters/http
Contains:
- handlers
- routes
- middleware
- validators

Handlers must:
- parse requests
- validate input
- call usecases
- return responses

Handlers must not contain business logic.

---

## adapters/postgres
Contains:
- GORM repositories
- persistence logic
- database queries

All database access must stay here.

---

## dto
Contains:
- request DTOs
- response DTOs

DTOs are not domain entities.

---

## mapper
Contains:
- DTO ↔ Domain transformations

---

# Dependency Injection

Dependencies must be initialized only in:

```text
/cmd/api/main.go
```

Do not instantiate repositories inside handlers.

---

# Database Rules

Database engine:
- PostgreSQL

ORM:
- GORM only

Use:
- explicit relations
- indexes
- transactions when needed

Avoid:
- raw SQL unless necessary
- hidden GORM magic
- automigration in production

---

# Naming Conventions

## Files

Use snake_case:

```text
crear_usuario.go
usuario_repository.go
jwt_middleware.go
```

---

## Interfaces

```go
type UsuarioRepository interface
```

Do not use:
- IUserRepository
- IUserService

---

## DTOs

Requests:
```go
CrearUsuarioRequest
```

Responses:
```go
UsuarioResponse
```

---

## UseCases

```go
CrearUsuarioUseCase
```

---

# API Rules

- REST API only
- JSON communication
- use proper HTTP status codes
- standardized API responses

Standard response format use from /internal/shared/response.go:

```json
{
  "success": true,
  "message": "OK",
  "data": {}
}
```

---

# Security

Required:
- JWT authentication
- bcrypt password hashing
- request validation
- environment variables for secrets

Never:
- hardcode secrets
- store plain text passwords

---

# Testing

Use:
- Testify

Requirements:
- unit tests for usecases
- repository tests isolated from business logic

Do not depend on real external services in unit tests.

---

# Logging

Use:
- Zap logger

Do not use:
```go
fmt.Println
```

for application logging.

---

# Migrations

Location:

```text
/migrations
```

Naming:

```text
000001_create_users.up.sql
000001_create_users.down.sql
```

---

# Development Commands

## Run API

```bash
go run ./cmd/api
```

---

## Run Tests

```bash
go test ./...
```

---

## Run Docker Environment

```bash
docker compose up --build
```

---

## Run Migrations

```bash
migrate -path migrations -database $DATABASE_URL up
```

---

# Module Creation Rules

Every new module must follow:

```text
/modules/module_name

    /domain
    /usecase
    /ports
    /adapters
    /dto
    /mapper
```

Do not create global business layers outside modules.

---

# Shared Code

Shared reusable components belong in:

```text
/internal/shared
```

Allowed shared components:
- config
- database
- logger
- middleware
- auth
- response helpers

Do not place business logic inside shared.

---

# Forbidden Patterns

Do not:
- place business logic in handlers
- access database directly from handlers
- couple modules directly
- create giant utility packages
- create god services
- bypass usecases

---

# Environment

Required files:

```text
.env
.env.example
```

---

# Docker

The project must run using:

```bash
docker compose up
```

---

# Priority Order

When making architectural decisions prioritize:

1. clarity
2. maintainability
3. simplicity
4. modularity
5. scalability
