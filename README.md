catatan_app_restful_API
RESTful API for notes management built with Go,  featuring clean layered architecture, repository pattern,  and dependency injection. Supports CRUD operations with archive functionality.

Architecture

This project follows a layered architecture pattern:
Handler → Service → Repository → Database

- Handler — HTTP request/response handling
- Service — Business logic
- Repository — Database access layer
- DTO — Request, response, and mapper
- Model — Domain struct

Tech Stack

- Language — Go (standard library)
- Database — MySQL
- Driver — go-sql-driver/mysql
- Config — joho/godotenv

Project Structure
notes-api/
├── server/
│   └── main.go
├── internal/
│   ├── apperror/
│   │   └── errors.go
│   ├── conect_db/
│   │   └── catatan_db_connection.go
│   ├── dto/
│   │   ├── catatan_mapper.go
│   │   ├── catatan_request.go
│   │   └── catatan_respons.go
│   ├── handler/
│   │   └── catatan_handler.go
│   ├── modul/
│   │   └── catatan_modul.go
│   ├── repository/
│   │   ├── interface.go
│   │   └── catatan_repo.go
│   ├── router/
│   │   └── catatan_route.go
│   └── service/
│       ├── interface.go
│       └── catatan_service.go
├── database/
│   └── catatan.sql
├── go.mod
└── go.sum

API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/catatan` | Create a new note |
| `GET` | `/catatan` | Get all active notes |
| `GET` | `/catatan?arsip=true` | Get all archived notes |
| `GET` | `/catatan/{id}` | Get note by ID |
| `PUT` | `/catatan/{id}` | Update note by ID |
| `DELETE` | `/catatan/{id}` | Delete note by ID |
| `PATCH` | `/catatan/{id}/arsip` | Archive a note |
| `PATCH` | `/catatan/{id}/unarsip` | Unarchive a note |

Design Patterns

- Layered Architecture — separation of concerns across handler, service, and repository
- Repository Pattern — database access abstraction via interfaces
- Dependency Injection — all dependencies injected through constructors in main.go
- Sentinel Error Pattern — centralized error definitions in apperror package
- DTO Pattern — separate request, response, and domain models

Features

- CRUD operations for notes
- Archive and unarchive notes
- Filter notes by archive status
- Graceful shutdown
- Database connection pooling
- Centralized error handling
- Interface-based design for testability
