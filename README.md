# Go URL Shortener 🔗

A high-performance URL shortener microservice built with **Go** (Golang), designed for scalability and ease of
deployment. Supports custom short URLs, analytics, and Redis caching.

![Go](https://img.shields.io/badge/Go-1.20%2B-blue)
![Redis](https://img.shields.io/badge/Redis-7.0%2B-red)
![License](https://img.shields.io/badge/License-MIT-green)

---

## Features ✨

- **Shorten long URLs** to concise, memorable links (e.g., `short.ly/abc123`).
- **Custom aliases**: Define your own short paths (e.g., `short.ly/mylink`).
- **Redis caching** for low-latency lookups.
- **REST API** with JSON responses.
- **Docker support** for easy deployment.
- Lightweight (<10MB Docker image).

## Quick Start 🚀

### Prerequisites

- Go 1.20+
- Redis (optional, for caching)

### Run Locally

```bash
# Clone the repo
git clone https://github.com/alexperezortuno/go-url-shortener.git
cd go-url-shortener
go mod tidy
```

# Start the server (default port :8080)

```bash
go run cmd/api/main.go
```

# With Redis (set REDIS_URL environment variable)

```bash
export REDIS_URL=redis://localhost:6379 && go run main.go
```

🚀 API Endpoints

| Endpoint                       | Method | Description                       | Example Request                                     | Example Response Body                                            |
|--------------------------------|--------|-----------------------------------|-----------------------------------------------------|------------------------------------------------------------------|
| `/url?short_url=<SHORT_CODE> ` | GET    | Get more data from URL            | -                                                   | {"long_url": "https://example.com", "short_url": "<SHORT_CODE>"} |
| `/url`                         | POST   | Create a short URL                | {"long_url": "https://example.com", "user_id": "1"} | {"short_url": "https://127.0.0.1/r/Eg4tQwFp"}                    |
| `/r/<SHORT_CODE>`              | GET    | Redirect to original URL          | -                                                   | redirect to url                                                  |
| `/health`                      | GET    | Get usage statistics (if enabled) |                                                     | {"message": "everything is ok"}                                  |

### Docker

#### Build the Docker image
```bash
docker build -t url-shortener:latest .
```

#### Create network

```bash
docker network create internal_net
```

#### Run the containers

```bash
docker-compose up --build -d
```

#### Stop the containers and re run

```bash
docker-compose down && docker-compose up --build -d
```
