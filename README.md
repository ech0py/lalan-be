# Lalan BE

A Go-based backend API for managing outdoor rental operations through an admin dashboard.  
Designed for scalability, maintainability, and clean architecture.

## Requirements

- Go 1.24.4 or higher  
- PostgreSQL  

## Getting Started

```bash
git clone https://github.com/braiyenmassora/lalan-be.git
cd lalan-be
go mod download
```

## Environment Configuration

Set up the `.env.dev` file before running the application:

```bash
# PostgreSQL database connection components (development)
DB_USER=<your_db_username>
DB_PASSWORD=<your_db_password>
DB_HOST=<your_db_host>
DB_PORT=<your_db_port>
DB_NAME=<your_db_name>

# JWT Secret Key (generate with: openssl rand -base64 32)
JWT_SECRET=<your_secret_key>

# Application Port
APP_PORT=8080
```

## Project Structure

```
lalan-be/
├── cmd/                    # Application entry point
├── internal/               # Core logic and modules
│   ├── config/             # App and database configuration
│   ├── handler/            # HTTP request handlers
│   ├── middleware/         # Authentication and middleware logic
│   ├── model/              # Data models
│   ├── repository/         # Database operations
│   ├── response/           # Response formatting utilities
│   ├── route/              # Route definitions
│   └── service/            # Business logic layer
├── migrations/             # Database migrations
├── pkg/                    # Shared helper packages
├── .env.dev                # Environment configuration (development)
├── go.mod                  # Go module definition
└── go.sum                  # Go module checksums
```

## Run Locally

### Install Air (Live Reload)

```bash
go install github.com/air-verse/air@latest
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### Run Application

```bash
# Run in development mode with live reload
make dev

# Or run manually
go run ./cmd/main.go

# Or build and execute binary
go build -o main ./cmd/main.go
./main
```

## Adding New Features

| Component  | Description                              | Location               |
|-------------|------------------------------------------|------------------------|
| Migration   | Manage database schema changes           | `migrations/`          |
| Model       | Define data structures                   | `internal/model/`      |
| Repository  | Implement database access logic          | `internal/repository/` |
| Service     | Handle business logic                    | `internal/service/`    |
| Handler     | Create HTTP request handlers             | `internal/handler/`    |
| Routes      | Register and manage API routes           | `internal/route/`      |
| Main        | Initialize and link all core components  | `cmd/main.go`          |

---

## Code Commenting Guidelines

```
{
  "initial": "Analyze this file to understand its role. Identify every const, var, type, func, method, init, and interface.",
  "after": "Remove all existing comments (// or /* */). No exceptions.",
  "new": "Reorder strictly using standard Go layout: package → imports → const → var → init → types (with methods) → interfaces → funcs (main last).\n\nBefore each element below, add exactly one block comment /* */ written in Indonesian:\n- Every const group\n- Every var group\n- Every type\n- Every function/method (exported or unexported if more than one line or has side effects)\n\nEach block comment must contain EXACTLY TWO LINES:\nLine 1: Tujuan utama fungsi/type/blok.\nLine 2: Hasil atau nilai yang dikembalikan.\n\nMandatory example format:\n\"\"\"\n/* Menginisialisasi koneksi ke PostgreSQL.\nMengembalikan konfigurasi database jika proses berhasil. */\n\"\"\"\n\nRestrictions:\n- Output comments must be in Indonesian only.\n- Maximum two lines.\n- Follow the exact two-sentence pattern with no creative variation.\n- Do not add or modify functional code.\n- Do not echo any code inside this instruction output."
}

```
