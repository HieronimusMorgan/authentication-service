# Authentication

## Description

**Authentication** is a robust and secure service designed to handle user authentication and authorization for modern web and mobile applications. This microservice provides features like user registration, login, token management, and role-based access control (RBAC) to ensure a secure and seamless authentication experience.

---

## Key Features

- **User Registration**: Allows users to create accounts with validation for secure data handling.
- **Login and Authentication**: Supports authentication using secure mechanisms like hashed passwords and JSON Web Tokens (JWT).
- **Token Management**:
  - Generate access tokens for user sessions.
  - Refresh tokens for session renewal without requiring re-login.
- **Role-Based Access Control (RBAC)**: Assign roles to users (e.g., admin, user) and manage permissions for protected resources.
- **Resource Management**: Enables defining and restricting access to application resources based on roles.
- **Internal Token Support**: Provides internal tokens for secure inter-service communication in microservices architecture.
- **Redis Integration**: Caches user data and tokens for optimized performance.
- **API Health Check**: Includes a default endpoint to check service health.

---

## Technology Stack

- **Backend Framework**: [Gin](https://gin-gonic.com/) - A high-performance HTTP web framework.
- **Database**: [PostgreSQL](https://www.postgresql.org/) or [SQLite](https://www.sqlite.org/) for persistence.
- **Authentication**: [JWT](https://jwt.io/) for stateless authentication.
- **Cache**: [Redis](https://redis.io/) for token storage.
- **ORM**: [GORM](https://gorm.io/) for database operations.

---

## Installation and Setup

### Prerequisites

- [Go](https://golang.org/doc/install) installed.
- A database setup (PostgreSQL or SQLite).
- Redis installed and running.

### Steps

1. Clone the repository:

   ```bash
   git clone https://github.com/HieronimusMorgan/authentication.git
   cd authentication
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure environment variables:

   - Create a `.env` file to store sensitive data like database connection strings, Redis connection, and JWT secret keys.

4. Run the application:

   ```bash
   go run main.go
   ```

---

## API Endpoints

### Public Routes

- `POST /auth/v1/register/internal-token`: Generate internal tokens for inter-service communication.

### Protected Routes (Require Authentication)

- `POST /auth/v1/register`: User registration.
- `POST /auth/v1/login`: User login.
- `POST /auth/v1/refresh`: Refresh token (to renew session).
- `GET /auth/v1/profile`: Fetch user profile (requires a valid token).

### Utility

- `GET /health`: Check service health.

---

## Internal Logic Overview

### 1. **Register**
- **Validates input** (username and password).
- **Hashes the password** for secure storage.
- **Creates a user** in the database.
- **Assigns resources** to the user.
- **Generates JWT tokens** for user sessions.
- **Caches user data** and tokens in Redis.

### 2. **Login**
- **Validates credentials** by comparing hashed passwords.
- **Generates JWT tokens** upon successful authentication.
- **Caches user session** data in Redis.

### 3. **Get Profile**
- **Extracts user claims** from the token.
- **Fetches user data** from the database using the ClientID.

### 4. **Register Internal Token**
- **Validates resource name** from the request.
- **Generates a unique internal token** for secure inter-service communication.
- **Stores the token** in the database.

---

## Redis Integration

### Purpose
- **Token Caching**: Speeds up token validation and reduces database load.
- **User Session Management**: Stores user session data for quick retrieval.

### Commands Used
- **Set Data**:
  ```go
  utils.SaveDataToRedis("key", "field", value)
  ```
- **Get Data**:
  ```go
  utils.GetDataFromRedis("key", "field")
  ```

---

## Contributions

Contributions are welcome! Please follow the guidelines in the `CONTRIBUTING.md` file to submit issues, suggest features, or create pull requests.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

