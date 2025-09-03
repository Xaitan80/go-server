# Go Server: Chirps API

A Go-based backend server for a "Chirps" application, similar to Twitter, allowing users to post short messages ("chirps"), manage users, and handle authentication. Includes admin and health endpoints.

---

## Features

- **Chirps**
  - Create a chirp: `POST /api/chirps`
  - List all chirps: `GET /api/chirps` with optional `author_id` filter and `sort` (`asc` or `desc`)
  - Retrieve a single chirp: `GET /api/chirps/{id}`
  - Delete a chirp: `DELETE /api/chirps/{id}`
- **Users**
  - Create a user: `POST /api/users`
  - Update a user: `PUT /api/users`
- **Authentication**
  - Login: `POST /api/login`
  - Refresh tokens: `POST /api/refresh`
  - Revoke tokens: `POST /api/revoke`
- **Admin & Metrics**
  - Get fileserver hits: `GET /admin/metrics`
  - Reset metrics: `POST /admin/reset`
- **Health Check**
  - Readiness endpoint: `GET /api/healthz`

---

## Dependencies

- Go ≥ 1.21  
- PostgreSQL  
- SQLC for type-safe SQL queries (`sqlc generate`)  
- JWT for authentication (`github.com/golang-jwt/jwt`)  
- Environment variables via `.env` (`github.com/joho/godotenv`)  
- PostgreSQL driver: `github.com/lib/pq`  

---

## Setup

1. **Clone repository**
```bash
git clone <repo-url>
cd go-server
