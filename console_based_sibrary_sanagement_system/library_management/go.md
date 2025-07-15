# Go Development Notes

This file contains Go-specific notes, conventions, and best practices for the Library Management System project.

## Go Coding Conventions
- Use `gofmt` or `go fmt` to format code consistently
- Follow idiomatic Go naming (CamelCase for exported, camelCase for unexported)
- Keep functions short and focused
- Use error handling idioms (`if err != nil { ... }`)
- Prefer composition over inheritance
- Use Go modules for dependency management

## Project Structure Tips
- `main.go`: Entry point, should be minimal (setup and start)
- `models/`: Define data structures (Book, Member, etc.)
- `services/`: Business logic and core operations
- `controllers/`: User interaction and command handling
- `docs/`: Documentation files

## Useful Go Commands
- Run the application:
  ```bash
  go run main.go
  ```
- Format code:
  ```bash
  go fmt ./...
  ```
- Run tests (if available):
  ```bash
  go test ./...
  ```
- Manage dependencies:
  ```bash
  go mod tidy
  ```

## Error Handling
- Always check errors returned by functions
- Return errors up the call stack with context if needed
- Use custom error types for domain-specific errors

## Further Reading
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Modules Reference](https://blog.golang.org/using-go-modules)

---
Feel free to add more Go-specific notes or tips as the project evolves.
