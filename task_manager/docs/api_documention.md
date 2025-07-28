# Task Manager API Documentation

## Table of Contents
1. [Base URL](#base-url)
2. [Authentication](#authentication)
3. [Endpoints](#endpoints)
   - [Get All Tasks](#1-get-all-tasks)
   - [Get Task by ID](#2-get-task-by-id)
   - [Add Task](#3-add-task)
   - [Update Task](#4-update-task)
   - [Remove Task](#5-remove-task)
   - [Register](#6-register)
   - [Login](#7-login)
   - [Promote User to Admin](#8-promote-user-to-admin)
4. [Error Response Example](#error-response-example)
5. [Rate Limiting](#rate-limiting)

---

## Base URL

```
http://localhost:8080/
```

---

## Authentication
- **Header:** `Authorization: Bearer <token>`
- **Description:** All endpoints (except `/register` and `/login`) require a valid JWT token for authentication. Include the token in the `Authorization` header of each request.

---

## Endpoints

### 1. Get All Tasks
- **Endpoint:** `GET /tasks`
- **Description:** Retrieve all tasks.
- **Query Parameters:**
  - `status` (optional): Filter tasks by status (e.g., `pending`, `in progress`, `completed`).
  - `due_date` (optional): Filter tasks by due date (e.g., `2025-08-01`).
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

### 6. Register
- **Endpoint:** `POST /register`
- **Description:** Register a new user.
- **Request Body:**
  ```json
  {
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "securepassword"
  }
  ```
- **Response:**
  ```json
  {
    "message": "User registered successfully",
    "user": {
      "id": "3",
      "username": "newuser",
      "email": "newuser@example.com",
      "role": "user"
    }
  }
  ```
- **Status Codes:**
  - 201 Created
  - 400 Bad Request

---

### 7. Login
- **Endpoint:** `POST /login`
- **Description:** Authenticate a user and receive a JWT.
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword"
  }
  ```
- **Response:**
  ```json
  {
    "message": "Login successful",
    "token": "jwt-token-here"
  }
  ```
- **Status Codes:**
  - 200 OK
  - 400 Bad Request

---

### 8. Promote User to Admin
- **Endpoint:** `POST /promote`
- **Description:** Promote a user to admin (admin only).
- **Request Body:**
  ```json
  {
    "user_id": "3"
  }
  ```
- **Response:**
  ```json
  {
    "message": "User promoted to admin successfully",
    "user": {
      "id": "3",
      "username": "newuser",
      "email": "newuser@example.com",
      "role": "admin"
    }
  }
  ```
- **Status Codes:**
  - 200 OK
  - 400 Bad Request
  - 404 Not Found

---

<!--
### (Not Implemented) Get All Users
### (Not Implemented) Get User by ID
### (Not Implemented) Update User
### (Not Implemented) Remove User
-->

## Error Response Example
```json
{
  "error": "Task not found"
}
```

---

## Rate Limiting
- **Description:** To prevent abuse, the API enforces rate limits.
- **Limits:**
  - 100 requests per minute per user.
- **Response for Exceeding Limit:**
  ```json
  {
    "error": "Rate limit exceeded. Try again later."
  }
  ```
