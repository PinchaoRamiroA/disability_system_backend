# SPECS.md

# Development Specifications

This document defines how features must be designed, implemented, tested, and integrated into the system.

All contributors and AI agents must follow these specifications.

---

# Feature Development Workflow

Every feature must follow this order:

1. Domain
2. DTOs
3. Ports
4. UseCase
5. Repository
6. HTTP Handler
7. Routes
8. Tests

Do not skip layers.

---

# Feature Structure

Each feature belongs to a module.

Example:

```text
/modules/incapacidades
```

Feature files should remain grouped by responsibility.

Example:

```text
/usecase/crear_incapacidad.go
/dto/crear_incapacidad_request.go
/dto/incapacidad_response.go
```

---

# UseCase Rules

Each usecase must:
- represent a single business action
- have one public Execute() method
- depend only on interfaces
- return controlled errors

Example:

```go
type CrearUsuarioUseCase struct {
    repo ports.UsuarioRepository
}

func (uc *CrearUsuarioUseCase) Execute(
    input dto.CrearUsuarioRequest,
) error {
}
```

---

# Handler Rules

Handlers must:
- parse request
- validate request
- call usecase
- return HTTP response

Handlers must not:
- contain business logic
- access database directly
- instantiate repositories

---

# Repository Rules

Repositories must:
- implement ports
- encapsulate persistence
- isolate GORM usage

Repositories must not:
- contain business rules
- perform HTTP operations

---

# DTO Rules

Use DTOs for:
- requests
- responses

Do not expose domain entities directly in API responses.

---

# Validation Rules

Validation must happen at:
- HTTP layer
- request DTOs

Use:
- Gin binding
- validator package if needed

---

# Error Handling

Use centralized application errors.

Standard error response:

```json
{
  "success": false,
  "message": "validation error",
  "error": {}
}
```

Avoid returning raw internal errors to clients.

---

# Database Standards

Use:
- UUID or BIGSERIAL consistently
- timestamps
- foreign keys
- indexes

Every table should include:
- created_at
- updated_at

Use soft delete only when necessary.

---

# GORM Standards

Allowed:
- explicit Preload()
- transactions
- scopes
- repository pattern

Avoid:
- deeply nested preload chains
- hidden automagic behavior
- business logic in models

---

# API Standards

API style:
- REST
- JSON only

Use proper status codes:
- 200 OK
- 201 Created
- 400 Bad Request
- 401 Unauthorized
- 404 Not Found
- 500 Internal Server Error

---

# Route Standards

Example:

```text
/api/v1/incapacidades
```

Use plural resource names.

---

# Authentication

Protected routes must use:
- JWT middleware

Passwords must:
- use bcrypt hashing

---

# Logging Standards

Use structured logging with Zap.

Every critical operation should log:
- request id
- module
- action
- error when applicable

---

# Testing Standards

Minimum required:
- unit tests for usecases

Recommended:
- repository integration tests
- handler tests

Use:
- Testify

---

# Dependency Injection

Dependency wiring must happen only in:

```text
/cmd/api/main.go
```

---

# Configuration

Environment variables must be loaded from:
- .env

Never hardcode:
- secrets
- database credentials
- JWT keys

---

# Migrations

Every schema change requires:
- up migration
- down migration

Migration names:

```text
000001_create_users.up.sql
000001_create_users.down.sql
```

---

# Documentation

Every module should include:
- README.md
- endpoint documentation

Recommended:
- Swagger/OpenAPI

---

# Pull Request Standards

Every PR should:
- compile successfully
- pass tests
- respect architecture
- avoid unrelated refactors

---

# Definition of Done

A feature is complete only if:

- architecture respected
- tests implemented
- migrations added if needed
- routes documented
- errors handled
- logs added when necessary
- no linting issues

---

# Performance Guidelines

Prefer:
- pagination
- indexed queries
- explicit selects

Avoid:
- N+1 queries
- loading unnecessary relations

---

# Security Guidelines

Never:
- expose internal errors
- trust client input
- store sensitive data unencrypted

Always:
- validate requests
- sanitize inputs
- protect private endpoints

---

# AI Agent Constraints

AI agents must:
- follow modular hexagonal architecture
- avoid creating unnecessary abstractions
- reuse existing patterns
- preserve module boundaries
- keep handlers thin
- keep usecases focused

When unsure:
- prioritize simplicity
- follow existing module patterns
