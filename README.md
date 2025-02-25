# 🔐 Authentication Service

## 📖 Description

**Authentication** is a robust and secure microservice designed to handle **user authentication and authorization** for modern **web and mobile applications**. This service provides features like **user registration, login, token management, and role-based access control (RBAC)** to ensure a seamless and secure authentication experience.

---

## ✨ Key Features

- ✅ **User Registration**: Allows users to create accounts with secure validation and password hashing.
- 🔑 **Login & Authentication**: Secure authentication using **hashed passwords** and **JSON Web Tokens (JWT)**.
- 🔄 **Token Management**:
  - Generate **access tokens** for user sessions.
  - Issue **refresh tokens** for session renewal without requiring re-login.
- 🎭 **Role-Based Access Control (RBAC)**: Assign roles (e.g., **admin, user**) and enforce permissions for protected resources.
- 📂 **Resource Management**: Define and restrict access to application resources **based on roles**.
- 🔒 **Internal Token Support**: Generate **internal tokens** for **secure inter-service communication**.
- 🚀 **Redis Integration**: Caches **user data** and **tokens** for optimized performance.
- ⚙️ **API Health Check**: Includes an endpoint to **check service health**.

---

## 🛠 Technology Stack

- **Backend Framework**: [Gin](https://gin-gonic.com/) - A high-performance HTTP web framework for Golang.
- **Database**: [PostgreSQL](https://www.postgresql.org/) / [SQLite](https://www.sqlite.org/)
- **Authentication**: [JWT](https://jwt.io/) for **stateless authentication**.
- **Cache**: [Redis](https://redis.io/) for **token storage** and caching.
- **ORM**: [GORM](https://gorm.io/) for **database operations**.

---

## 📦 Installation and Setup

### Prerequisites

- Install **[Go](https://golang.org/doc/install)**.
- Set up **PostgreSQL or SQLite**.
- Install and run **Redis**.

### Steps to Run

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
   - Create a `.env` file to store **database connection strings, Redis connection, and JWT secrets**.

4. Run the application:
   ```bash
   go run main.go
   ```

---

## 🔗 API Endpoints

### 🔓 Public Routes
- `POST /v1/register/internal-token` → Generate **internal tokens** for inter-service communication.

### 🔒 Protected Routes (Require Authentication)
- `POST /v1/register` → **User registration**.
- `POST /v1/login` → **User login**.
- `POST /v1/refresh` → **Refresh token** to extend session.
- `GET /v1/profile` → **Fetch user profile** (requires valid token).

### ⚙️ Utility
- `GET /health` → **Service health check**.

---

## 🔍 Internal Logic Overview

### 1️⃣ **User Registration**
- **Validates input** (username & password).
- **Hashes the password** before storing.
- **Creates user** in the database.
- **Assigns roles & resources** to the user.
- **Generates JWT tokens**.
- **Caches user data & tokens** in Redis.

### 2️⃣ **User Login**
- **Validates credentials** by comparing hashed passwords.
- **Generates JWT tokens** on successful authentication.
- **Stores session data** in Redis.

### 3️⃣ **Fetching User Profile**
- **Extracts user claims** from token.
- **Fetches user data** using ClientID.

### 4️⃣ **Internal Token Registration**
- **Validates resource name**.
- **Generates unique internal token**.
- **Stores token** securely in database.

---

## ⚡ Redis Integration

### 🚀 Purpose
- **Token Caching** → Speeds up validation & reduces DB queries.
- **Session Management** → Stores session data for quick retrieval.

### 🛠 Commands Used

- **Save Data to Redis**:
  ```go
  utils.SaveDataToRedis("key", "field", value)
  ```

- **Retrieve Data from Redis**:
  ```go
  var token string
  utils.GetDataFromRedis("key", "field", &token)
  ```

---

## 🤝 Contributions

Contributions are **welcome**! Follow these steps:

1. Fork the repository.
2. Create a new branch.
3. Submit a pull request with **detailed changes**.

For major updates, **open an issue** first to discuss your proposal.

---

## 📜 License

This project is licensed under the **MIT License**. See the `LICENSE` file for more details.

---

This **Authentication Service** is designed to **enhance security and efficiency** in modern applications, ensuring seamless user authentication and authorization. 🚀🔐
