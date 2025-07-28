# Task Manager

A simple task and user management API built with Go, Gin, and MongoDB.

## Features
- User registration and authentication (JWT)
- Admin user promotion
- Task CRUD operations (create, read, update, delete)
- Role-based access control (admin/user)
- RESTful API design

## Tech Stack
- Go (Golang)
- Gin (HTTP web framework)
- MongoDB (database)

## Getting Started

### Prerequisites
- Go 1.18+
- MongoDB instance (local or cloud)

### Installation & Setup
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd task_manager
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your environment variables (MongoDB URI, JWT secret, etc.) as needed.
4. Run the application:
   ```bash
   go run main.go
   ```

### API Usage
- The API runs by default at `http://localhost:8080/`.
- See `docs/api_documention.md` for full API details and endpoints.

## Folder Structure
- `Domain/` — Domain models and interfaces
- `Repository/` — Database access logic
- `Usecases/` — Business logic
- `Infrastructure/` — Services (JWT, password, middleware)
- `Delivery/` — HTTP handlers, controllers, routers
- `docs/` — API documentation

## License
MIT (or specify your license)
