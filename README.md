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
- **API Health Check**: Includes a default endpoint to check service health.

---

## Technology Stack

- **Backend Framework**: [Gin](https://gin-gonic.com/) - A high-performance HTTP web framework.
- **Database**: [PostgreSQL](https://www.postgresql.org/) or [SQLite](https://www.sqlite.org/) for persistence.
- **Authentication**: [JWT](https://jwt.io/) for stateless authentication.
- **Cache**: [Redis](https://redis.io/) for token storage (optional).
- **ORM**: [GORM](https://gorm.io/) for database operations.

---

## Installation and Setup

### Prerequisites

- [Go](https://golang.org/doc/install) installed.
- A database setup (PostgreSQL or SQLite).

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/<your-username>/authentication.git
   cd authentication
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure environment variables:
   - Create a `.env` file to store sensitive data like database connection strings and JWT secret keys.

4. Run the application:
   ```bash
   go run main.go
   ```

---

## API Endpoints

### Public Routes

- `POST /auth/v1/register`: User registration.
- `POST /auth/v1/login`: User login.
- `POST /auth/v1/refresh`: Token refresh.

### Protected Routes

- `GET /auth/v1/profile`: Fetch user profile (requires a valid token).
- `POST /auth/v1/register/internal-token`: Generate internal tokens for service communication.

### Utility

- `GET /health`: Check service health.

---

## Contributions

Contributions are welcome! Please follow the guidelines in the `CONTRIBUTING.md` file to submit issues, suggest features, or create pull requests.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

