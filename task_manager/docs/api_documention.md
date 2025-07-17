# Task Manager API Documentation

## Base URL

```
http://localhost:8080/
```

---

## Endpoints

### 1. Get All Tasks
- **Endpoint:** `GET /tasks`
- **Description:** Retrieve all tasks.
- **Response:**
```json
{
  "tasks": [
    {
      "id": "1",
      "title": "Write code",
      "description": "Finish Go project",
      "due_date": "2025-08-01T00:00:00Z",
      "status": "pending"
    }
  ]
}
```
- **Status Codes:**
  - 200 OK

---

### 2. Get Task by ID
- **Endpoint:** `GET /tasks/:id`
- **Description:** Retrieve a single task by its ID.
- **Response:**
```json
{
  "task": {
    "id": "1",
    "title": "Write code",
    "description": "Finish Go project",
    "due_date": "2025-08-01T00:00:00Z",
    "status": "pending"
  }
}
```
- **Status Codes:**
  - 200 OK
  - 404 Not Found

---

### 3. Add Task
- **Endpoint:** `POST /tasks`
- **Description:** Add a new task.
- **Request Body:**
```json
{
  "id": "1",
  "title": "Write code",
  "description": "Finish Go project",
  "due_date": "2025-08-01T00:00:00Z",
  "status": "pending"
}
```
- **Response:**
```json
{
  "message": "Task added successfully",
  "task": {
    "id": "1",
    "title": "Write code",
    "description": "Finish Go project",
    "due_date": "2025-08-01T00:00:00Z",
    "status": "pending"
  }
}
```
- **Status Codes:**
  - 201 Created
  - 400 Bad Request

---

### 4. Update Task
- **Endpoint:** `PUT /tasks/:id`
- **Description:** Update an existing task.
- **Request Body:**
```json
{
  "title": "Write more code",
  "description": "Finish Go project and tests",
  "due_date": "2025-08-02T00:00:00Z",
  "status": "in progress"
}
```
- **Response:**
```json
{
  "message": "Task updated successfully",
  "task": {
    "id": "1",
    "title": "Write more code",
    "description": "Finish Go project and tests",
    "due_date": "2025-08-02T00:00:00Z",
    "status": "in progress"
  }
}
```
- **Status Codes:**
  - 200 OK
  - 400 Bad Request
  - 404 Not Found

---

### 5. Remove Task
- **Endpoint:** `DELETE /tasks/:id`
- **Description:** Remove a task by its ID.
- **Response:**
```json
{
  "message": "Task removed successfully"
}
```
- **Status Codes:**
  - 200 OK
  - 404 Not Found

---

## Error Response Example
```json
{
  "error": "Task not found"
}
```
