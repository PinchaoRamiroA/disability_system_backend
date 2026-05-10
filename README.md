# Backend - Sistema de Gestión de Incapacidades

Backend empresarial desarrollado en Go siguiendo principios de Arquitectura Hexagonal (Ports & Adapters), modularidad y buenas prácticas de ingeniería de software.

---

# 📌 Objetivo del Proyecto

Este sistema tiene como objetivo gestionar de manera centralizada:

- Incapacidades médicas
- Transcripción y validación documental
- Gestión de pagos y cobros
- Seguimiento administrativo
- Historial de eventos
- Notificaciones
- Gestión de usuarios y roles

Diseñado para ser:

- Escalable
- Mantenible
- Modular
- Seguro
- Performante
- Preparado para crecimiento futuro

---

# 🏛️ Arquitectura

El proyecto implementa:

## ✅ Arquitectura Hexagonal (Ports & Adapters)

Separando claramente:

| Capa | Responsabilidad |
|---|---|
| Domain | Reglas de negocio |
| UseCases | Casos de uso |
| Ports | Contratos/interfaces |
| Adapters | Infraestructura |
| Shared | Componentes reutilizables |

---

# 📦 Estructura del Proyecto

```text
/disability_system_backend
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   │
│   ├── shared/
│   │   ├── auth/
│   │   ├── config/
│   │   ├── database/
│   │   ├── errors/
│   │   ├── logger/
│   │   ├── middleware/
│   │   ├── response/
│   │   └── utils/
│   │
│   ├── modules/
│   │
│   │   ├── auth/
│   │   │   ├── domain/
│   │   │   ├── usecase/
│   │   │   ├── ports/
│   │   │   ├── adapters/
│   │   │   │   ├── http/
│   │   │   │   └── postgres/
│   │   │   ├── dto/
│   │   │   └── mapper/
│   │   │
│   │   ├── usuarios/
│   │   ├── incapacidades/
│   │   ├── cobros/
│   │   ├── historial/
│   │   └── notificaciones/
│
├── migrations/
│
├── scripts/
│
├── tests/
│
├── deployments/
│
├── docs/
│
├── .env
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── README.md
```

---

# 🧠 Principios Arquitectónicos

## ✅ Arquitectura Modular

Cada módulo es independiente y encapsula:

- dominio
- persistencia
- casos de uso
- contratos
- adaptadores

---

## ✅ Separación de Responsabilidades

El dominio NO conoce:

- GORM
- PostgreSQL
- HTTP
- Frameworks
- Infraestructura

---

## ✅ Relaciones Cross-Module

Los modelos GORM NO deben tener relaciones entre módulos.

### ✔ Permitido

```text
usuarios
 ├── Usuario
 ├── Rol
 └── Empleado
```

### ❌ Evitar

```text
Incapacidad -> UsuarioModel
```

En su lugar:

```go
IDUsuario uint64
```

Y resolver mediante:
- repositories
- usecases
- joins explícitos
- DTOs

---

# 🗄️ Base de Datos

## Motor

- PostgreSQL

## ORM

- GORM

## Migraciones

- golang-migrate

---

# ⚠️ Importante

El proyecto NO utiliza:

```go
AutoMigrate()
```

como fuente oficial del schema.

---

## Fuente Oficial del Schema

```text
/migrations
```

---

# 📁 Organización de Persistencia

## Domain

Entidades puras:

```go
type Usuario struct {
    ID uint64
}
```

---

## Models GORM

Infraestructura:

```go
type UsuarioModel struct {
    IDUsuario uint64 `gorm:"primaryKey"`
}
```

Ubicación:

```text
/adapters/postgres/models
```

---

# 🚀 Tecnologías Utilizadas

| Tecnología | Uso |
|---|---|
| Go | Backend principal |
| PostgreSQL | Base de datos |
| GORM | ORM |
| golang-migrate | Migraciones |
| Docker | Contenedores |
| JWT | Autenticación |
| Gin | HTTP Framework |
| Makefile | Automatización |
| Swagger | Documentación API |

---

# 🔐 Autenticación

El sistema utiliza:

- JWT Access Tokens
- Middleware de autorización
- Roles y permisos
- Guards por módulo

---

# 👥 Módulos del Sistema

| Módulo | Responsabilidad |
|---|---|
| auth | autenticación/autorización |
| usuarios | usuarios y roles |
| incapacidades | gestión de incapacidades |
| cobros | pagos y seguimiento |
| historial | auditoría e historial |
| notificaciones | notificaciones del sistema |

---

# 🧾 Reglas de Negocio

## Usuarios

- Todo usuario debe tener un rol
- El correo debe ser único
- El documento debe ser único

---

## Incapacidades

- Toda incapacidad pertenece a un usuario
- Toda incapacidad tiene estado
- La fecha final no puede ser menor a la inicial
- Puede requerir documentos obligatorios
- Puede tener seguimiento y pagos asociados

---

## Documentos

- Deben estar asociados a una incapacidad
- Pueden ser validados por Gestión Humana
- Deben registrar estado y formato

---

## Pagos

- El valor debe ser mayor o igual a 0
- Deben estar asociados a una entidad
- Pueden ser conciliados

---

# ⚙️ Variables de Entorno

## `.env`

```env
APP_PORT=PORT

DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=password
DB_NAME=disability_system_db
DB_SSLMODE=disable

DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=5
DATABASE_URL=postgres://root:password@localhost:5432/disability_system_db?sslmode=disable

JWT_SECRET=secret
JWT_EXPIRATION_HOURS=24
```

---

# 🐳 Docker

## Levantar servicios

```bash
docker-compose up -d
```

---

# 📦 Instalación Local

## 1️⃣ Clonar repositorio

```bash
git clone <repository>
```

---

## 2️⃣ Entrar al proyecto

```bash
cd disability_system_backend
```

---

## 3️⃣ Instalar dependencias

```bash
go mod tidy
```

---

## 4️⃣ Configurar variables de entorno

```bash
cp .env.example .env
```

---

## 5️⃣ Ejecutar migraciones

```bash
make migrate-up
```

---

## 6️⃣ Ejecutar proyecto

```bash
go run cmd/api/main.go
```

---

# 🛠️ Comandos Makefile

```bash
make run
make dev
make test
make lint
make migrate-up
make migrate-down
make migrate-force
```

---

# 🧪 Testing

El proyecto debe incluir:

- Unit tests
- Integration tests
- Repository tests
- UseCase tests

---

# 📚 Convenciones del Proyecto

## Naming

### Models

```go
UsuarioModel
```

---

### Domain

```go
Usuario
```

---

### DTOs

```go
CreateUsuarioRequest
UsuarioResponse
```

---

# 🚫 Reglas Importantes

## ❌ No usar AutoMigrate en producción

---

## ❌ No exponer modelos GORM en APIs

Siempre usar DTOs.

---

## ❌ No usar lógica de negocio en handlers HTTP

---

## ❌ No usar relaciones GORM cross-module

---

## ✅ Mantener módulos desacoplados

---

## ✅ Mantener dominio puro

---

# 📈 Escalabilidad

La arquitectura fue diseñada para permitir:

- Modular Monolith
- Evolución futura a microservicios
- Horizontal scaling
- Nuevos módulos
- Nuevos adaptadores

---

# 🔍 Health Check

```http
GET /health
```

---

# 📖 Documentación API

Swagger/OpenAPI será generado automáticamente.

---

# 🧠 Filosofía del Proyecto

Este backend prioriza:

- claridad
- mantenibilidad
- separación de responsabilidades
- performance
- consistencia
- evolución futura

---

# 👨‍💻 Estándares de Desarrollo

- Clean Code
- SOLID
- Hexagonal Architecture
- DDD-inspired modularity
- Repository Pattern
- DTO Pattern
- Explicit Migrations
- Dependency Injection

---

# 📌 Estado del Proyecto

🚧 En desarrollo activo.

---

# 📄 Licencia

Privado / Uso interno.