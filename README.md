# Go Role-Based Authorization API

This project is a production-ready REST API built in Go that demonstrates a secure, role-based authorization system. The application features two distinct roles: `admin` and `user`. Admins have access to a special admin dashboard, while regular users can only access a standard user dashboard.

The application is built with a modern, layered architecture and follows professional Go development standards, including dependency injection, structured logging, and secure credential handling.

## Features

-   **Role-Based Access Control (RBAC):** Secure endpoints using JWTs and role-based middleware.
-   **Secure Authentication:** Uses JWTs for stateless authentication and `bcrypt` for password hashing.
-   **CLI for Admin Creation:** A secure, command-line interface for creating the initial admin user, preventing exposure of admin creation logic via a public API endpoint.
-   **Layered Architecture:** Clear separation of concerns between API handlers, business logic (services), and data access (repositories).
-   **Structured Logging:** Uses Go's standard `slog` library for JSON-formatted, structured logs.
-   **Configuration Management:** Loads configuration from a `.env` file for easy setup across different environments.
-   **Database:** Uses SQLite with GORM for simple setup and powerful ORM features.
-   **Robust Routing:** Utilizes `chi` for flexible and powerful HTTP routing and middleware management.
-   **Input Validation:** Automatic request body validation using `go-playground/validator`.

## Project Structure

```text
go-auth-app/
├── cmd/api/main.go         # Application entrypoint (web server & CLI)
├── internal/               # All private application logic
│   ├── api/                # HTTP handlers, routing, and middleware
│   ├── auth/               # JWT generation and password hashing
│   ├── database/           # Database connection and migration
│   ├── model/              # Data models (structs)
│   ├── repository/         # Data access layer
│   └── service/            # Business logic layer
├── .env.example            # Example environment variables
├── go.mod                  # Go module dependencies
└── README.md
```



# Getting Started

## Prerequisites

* Go 1.21 or later
* (Optional but recommended) `curl` and `jq` for command-line testing

---

## 1. Clone the Repository

```bash
git clone https://github.com/cliffdoyle/go-auth-app
cd go-auth-app
```

---

## 2. Configure Environment

Copy the example environment file and update it if necessary. The default values are fine for local testing.

```bash
cp .env.example .env
```

---

## 3. Install Dependencies

Go modules will automatically handle dependency installation when you build or run the application. You can also fetch them manually:

```bash
go mod tidy
```

---

# How to Run and Test the Application

The application has two modes:

* CLI mode for administrative tasks
* Web server mode for running the API

---

## Step 1: Create the Admin User (CLI Mode)

Create the initial administrative user using the secure CLI command:

```bash
go run . create-admin --name="Super Admin" --email="admin@example.com" --password="a-very-secure-password"
```

You should see a success message in your terminal confirming that the admin user was created.

---

## Step 2: Start the Web Server

Run the application without any commands to start the API server:

```bash
go run .
```

The server will start on:
[http://localhost:8080](http://localhost:8080)

---

## Step 3: Test the API Endpoints

You can use `curl` or Postman to test the API.

---

### 1. Register a Regular User

This user will have the `user` role:

```bash
curl -X POST http://localhost:8080/api/register \
-H "Content-Type: application/json" \
-d '{
    "name": "Normal User",
    "email": "user@example.com",
    "password": "password123"
}'
```

---

### 2. Log In and Get Tokens

#### Log in as Admin:

```bash
ADMIN_TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{"email": "admin@example.com", "password": "a-very-secure-password"}' | jq -r .token)
```

#### Log in as User:

```bash
USER_TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
-H "Content-Type: application/json" \
-d '{"email": "user@example.com", "password": "password123"}' | jq -r .token)
```

---

### 3. Test Dashboard Access

#### Admin accessing the Admin Dashboard (SUCCESS)

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8080/api/dashboard/admin
```

**Expected Response:**

```json
{
  "admin_email": "admin@example.com",
  "message": "Welcome to the Admin Dashboard!"
}
```

---

#### User accessing the User Dashboard (SUCCESS)

```bash
curl -H "Authorization: Bearer $USER_TOKEN" http://localhost:8080/api/dashboard/user
```

**Expected Response:**

```json
{
  "email": "user@example.com",
  "message": "Welcome to your user dashboard!"
}
```

---

#### User trying to access the Admin Dashboard (FAILURE - 403 Forbidden)

This is the key test for role-based authorization.

```bash
curl -i -H "Authorization: Bearer $USER_TOKEN" http://localhost:8080/api/dashboard/admin
```

**Expected Response:**

* HTTP/1.1 403 Forbidden
* Body:

```json
{
  "error": "Forbidden"
}
```

---

## Notes

* The `.env.example` file contains default values for development.
* Be sure not to commit your `.env` file to version control.
* Use `jq` to cleanly extract tokens from login responses in `curl`.

---

## License

MIT License. See the `LICENSE` file for details.
