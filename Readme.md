# BalkanID File Vault

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?style=for-the-badge&logo=react)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-4169E1?style=for-the-badge&logo=postgresql)
![Docker](https://img.shields.io/badge/Docker-20.10-2496ED?style=for-the-badge&logo=docker)
![TailwindCSS](https://img.shields.io/badge/Tailwind_CSS-3-38B2AC?style=for-the-badge&logo=tailwind-css)

A production-grade, secure file vault system built with Go and React. This application supports efficient, deduplicated file storage, a powerful search and filtering API, and a granular, role-based access control (RBAC) system for secure file sharing and management.

## Features

### Core User Features

- [x] **Secure User Authentication**: JWT-based authentication with password hashing.
- [x] **Deduplicated File Storage**: Content-based hashing (SHA-256) to prevent duplicate data storage, saving space.
- [x] **Multi-File Uploads**: Supports single and multiple file uploads with a drag-and-drop UI.
- [x] **Rich File Management**:
  - [x] List, preview, and download files.
  - [x] Grid and List view options.
  - [ ] Manage file tags (add/remove).
- [x] **Advanced Sharing Controls**:
  - [x] Create public, shareable links.
  - [x] Share files with specific users by email.
  - [x] View and revoke shares.
- [ ] **Powerful Search**: Debounced, multi-field search (filename, tags, date) with database-level optimizations.
- [x] **Storage Statistics**: Users can view their storage usage, including savings from deduplication.
- [x] **Light/Dark Mode**: A theme toggle for user comfort.

### Admin Features

- [x] **Role-Based Access Control (RBAC)**: A full-featured RBAC system protecting all sensitive routes.
- [x] **System-Wide Dashboard**: Admins can view all files, system-wide statistics, and user information.
- [x] **User Management**: Admins can configure user-specific storage quotas.
- [x] **Audit Logging**: All critical actions (uploads, deletes, shares) are logged for security and compliance.

## 🛠️ Tech Stack

| Category     | Technology                            |
| ------------ | ------------------------------------- |
| **Backend**  | Go (Golang)                           |
| **Frontend** | React, TypeScript, Vite, Tailwind CSS |
| **Database** | PostgreSQL                            |
| **Storage**  | MinIO (S3-Compatible Object Storage)  |
| **API**      | REST                                  |
| **DevOps**   | Docker, Docker Compose, Makefile      |

## Getting Started

### Prerequisites

- Git
- Docker and Docker Compose

### Setup & Installation

1.  **Clone the repository:**

    ```bash
    git clone <your-repo-url>
    cd file-vault
    ```

2.  **Configure Environment Variables:**
    Copy the example environment file and customize if needed. The default values are configured to work with Docker Compose out of the box.

    ```bash
    cp .env.example .env
    ```

3.  **Build and Run the Application:**
    First, build all services in production mode:

    ```bash
    make build
    ```

    Then, apply database migrations:

    ```bash
    make migrate-up
    ```

    Finally, seed the database with initial roles and permissions:

    ```bash
    make seed
    ```

4.  **Alternative: Development Mode (Optional):**
    For development with hot-reloading, you can use:
    ```bash
    make up
    ```
    Then run `make seed` after the containers are running.

### Accessing the Application

- **Frontend (React App)**: [http://localhost:3000](http://localhost:3000)
- **Backend (Go API)**: [http://localhost:8080](http://localhost:8080)
- **MinIO Console (Object Storage UI)**: [http://localhost:9001](http://localhost:9001) (Use credentials from your `.env` file).

### Makefile Commands

This project uses a `Makefile` to simplify common development tasks.

| Command             | Description                                                           |
| ------------------- | --------------------------------------------------------------------- |
| `make up`           | Starts all services in development mode with hot-reloading.           |
| `make down`         | Stops and removes all containers, networks, and volumes.              |
| `make build`        | Builds and starts all services in production mode (no hot-reloading). |
| `make logs`         | Tails the logs of all running services.                               |
| `make seed`         | Seeds the database with initial roles and permissions.                |
| `make sqlc`         | Regenerates Go code from your SQL queries.                            |
| `make migrate-up`   | Applies all database migrations.                                      |
| `make migrate-down` | Rolls back all database migrations.                                   |
