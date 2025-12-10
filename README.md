# 🎫 Ticket Booking Service

**Ticket Booking Service** is a robust backend application designed for handling high-traffic event ticketing systems ("Ticket War"). It handles user authentication, event management, order processing with race-condition handling (Redis), and PDF ticket generation.

---

## 🚀 Tech Stack

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Language** | ![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go&logoColor=white) | Core programming language |
| **Framework** | ![Fiber](https://img.shields.io/badge/Fiber-v2-000000?style=flat) | High performance web framework |
| **Database** | ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql&logoColor=white) | Primary relational database |
| **Caching** | ![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat&logo=redis&logoColor=white) | Caching and Atomic Locks for stock management |
| **Storage** | ![MinIO](https://img.shields.io/badge/MinIO-S3-C72C48?style=flat&logo=minio&logoColor=white) | S3-compatible object storage for images/files |
| **Broker** | ![RabbitMQ](https://img.shields.io/badge/RabbitMQ-3.12-FF6600?style=flat&logo=rabbitmq&logoColor=white) | Asynchronous messaging (queue) |
| **Container** | ![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker&logoColor=white) | Containerization for easy deployment |

---

## ✨ Features

### 🔐 Authentication & Users
- **Register & Login** (JWT Based)
- **Profile Management** (View Profile)
- **Avatar Upload** (Stored in MinIO/S3)

### 📅 Event Management
- **Create Events** with images
- **Browse Events** (List & Detail view)
- **Stock Management** (Real-time availability)

### 🛒 Ordering System (The "War" Part)
- **High Concurrency Order Handling**: Uses Redlock/Redis atomic operations to prevent overselling ("race conditions").
- **Booking Flow**: Reserve ticket -> Payment Webhook -> Confirm.
- **Ticket Generation**: Generates PDF tickets with unique QR/Barcodes.

---

## 🛠️ How to Run

### prerequisites
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running.

### 1. Clone & Configure
Copy the example environment file:
```bash
cp .env.example .env
```
> **Note**: The default `.env.example` is configured to work out-of-the-box with Docker.

### 2. Start Services
Run the entire stack (App, DB, Redis, MinIO, RabbitMQ) with one command:
```bash
docker-compose up -d --build
```

### 3. Access
- **API Server**: `http://localhost:8000`
- **MinIO Console**: `http://localhost:9001` (User/Pass: `minioadmin` by default if not changed)
- **RabbitMQ Management**: `http://localhost:15672` (User/Pass: `guest`)

---

## 📂 Project Structure
```
├── 📂 cmd          # Entrypoint (main.go)
├── 📂 configs      # Configuration loader
├── 📂 internal
│   ├── 📂 app      # Server setup & Routing
│   ├── 📂 domain   # Database Models & Entities
│   ├── 📂 features # Business Logic (Usecases, Handlers)
│   └── 📂 platform # Infra setup (DB, Redis, MinIO, etc.)
├── 📄 Dockerfile
└── 📄 docker-compose.yml
```

---

Made with ❤️ by developer.
