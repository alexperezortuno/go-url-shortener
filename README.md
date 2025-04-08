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

| Endpoint        | Method | Description                       | Example Request Body           |
|-----------------|--------|-----------------------------------|--------------------------------|
| `/shorten`      | POST   | Create a short URL                | {"url": "https://example.com"} |
| `/{short_code}` | GET    | Redirect to original URL          | -                              |
| `/analytics`    | GET    | Get usage statistics (if enabled) | -                              |
