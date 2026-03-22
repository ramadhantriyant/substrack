# Substrack

A self-hosted subscription tracker with a Go REST API backend and a Svelte web UI.

## Features

- JWT-based authentication with refresh tokens
- Track subscriptions with billing cycles, costs, and statuses
- Organize subscriptions by categories
- Per-user subscription and category management
- Argon2id password hashing
- Graceful shutdown
- Compiled Svelte UI served directly from the Go binary

## Tech Stack

**Backend**
- **Go** 1.26 (stdlib `net/http`)
- **SQLite** via `mattn/go-sqlite3` (CGO)
- **sqlc** for type-safe query generation
- **goose** for database migrations
- **argon2id** (`alexedwards/argon2id`) for password hashing
- **JWT** (`golang-jwt/jwt/v5`, HS256) for access tokens

**Frontend**
- **Svelte 5** with TypeScript
- **Vite** for bundling
- Embedded into the Go binary via `//go:embed`

## Getting Started

### Prerequisites

- Go 1.26+
- C compiler (required for CGO/SQLite)
- [Bun](https://bun.sh) (for building the UI)

### Run Locally

```bash
# 1. Build the UI
cd ui && bun install && bun run build && cd ..

# 2. Run the server
export JWT_SECRET=$(openssl rand -base64 32)
go run .
```

The server listens on `:8080`. Open `http://localhost:8080` in your browser.
The SQLite database is created at `data/substrack.db` and migrations run automatically on startup.

### Run with Docker

```bash
docker build -t substrack .
docker run -p 8080:8080 -e JWT_SECRET=your_secret -v substrack_data:/app/data substrack
```

The Docker build compiles the Svelte UI and embeds it into the Go binary automatically.

## Environment Variables

| Variable     | Required | Description                                             |
|--------------|----------|---------------------------------------------------------|
| `JWT_SECRET` | Yes      | Secret key for signing JWTs (min 32 bytes recommended)  |

## API Reference

All endpoints except `/auth/register`, `/auth/login`, and `/auth/refresh` require a valid `Authorization: Bearer <token>` header.

### Auth

| Method | Path             | Description                           |
|--------|------------------|---------------------------------------|
| POST   | `/auth/register` | Register a new user                   |
| POST   | `/auth/login`    | Login, returns access + refresh token |
| POST   | `/auth/refresh`  | Issue a new access token              |
| POST   | `/auth/logout`   | Revoke refresh token                  |

#### Register
```json
POST /auth/register
{
  "email": "user@example.com",
  "name": "Alice",
  "password": "password123"
}
```

#### Login
```json
POST /auth/login
{
  "email": "user@example.com",
  "password": "password123"
}
```
Response:
```json
{
  "access_token": "eyJ...",
  "refresh_token": "abc123...",
  "token_type": "Bearer"
}
```

Access tokens expire in **1 hour**. Refresh tokens expire in **24 hours**.

#### Refresh
```json
POST /auth/refresh
{
  "refresh_token": "abc123..."
}
```

#### Logout
```json
POST /auth/logout
Authorization: Bearer <token>
{
  "refresh_token": "abc123..."
}
```

---

### User

| Method | Path                    | Description              |
|--------|-------------------------|--------------------------|
| GET    | `/api/user/me`          | Get current user profile |
| PUT    | `/api/user/me`          | Update profile           |
| PUT    | `/api/user/me/password` | Change password          |
| DELETE | `/api/user/me`          | Delete account           |

---

### User Subscriptions

| Method | Path                                | Description                             |
|--------|-------------------------------------|-----------------------------------------|
| GET    | `/api/user/me/subscription`         | List subscriptions for current user     |
| POST   | `/api/user/me/subscription/{id}`    | Add a subscription to current user      |
| DELETE | `/api/user/me/subscription/{id}`    | Remove a subscription from current user |

---

### User Categories

| Method | Path                          | Description                         |
|--------|-------------------------------|-------------------------------------|
| GET    | `/api/user/me/category`       | List categories for current user    |
| POST   | `/api/user/me/category/{id}`  | Add a category to current user      |
| DELETE | `/api/user/me/category/{id}`  | Remove a category from current user |

---

### Categories

| Method | Path                             | Description                |
|--------|----------------------------------|----------------------------|
| GET    | `/api/category`                  | List all categories        |
| GET    | `/api/category/id/{id}`          | Get category by ID         |
| GET    | `/api/category/name/{name}`      | Get category by name       |
| POST   | `/api/category`                  | Create a category          |
| PUT    | `/api/category/{id}`             | Update a category          |
| PUT    | `/api/category/{id}/name`        | Update category name       |
| PUT    | `/api/category/{id}/description` | Update category description|
| DELETE | `/api/category/{id}`             | Delete a category          |

---

### Subscriptions

| Method | Path                                  | Description                                       |
|--------|---------------------------------------|---------------------------------------------------|
| GET    | `/api/subscription`                   | List all subscriptions (optional `?category_id=`) |
| GET    | `/api/subscription/active`            | List active subscriptions                         |
| GET    | `/api/subscription/expired`           | List expired subscriptions                        |
| GET    | `/api/subscription/cycle/{billCycle}` | List by billing cycle                             |
| GET    | `/api/subscription/{id}`              | Get subscription by ID                            |
| POST   | `/api/subscription`                   | Create a subscription                             |
| PUT    | `/api/subscription/{id}`              | Update a subscription                             |
| PUT    | `/api/subscription/{id}/status`       | Update subscription status                        |
| PUT    | `/api/subscription/{id}/cost`         | Update subscription cost                          |
| PATCH  | `/api/subscription/{id}/pause`        | Pause a subscription                              |
| DELETE | `/api/subscription/{id}`              | Delete a subscription                             |

## Project Structure

```
.
├── main.go                        # Entry point, DB init, server start
├── server.go                      # Route registration, static UI serving
├── Dockerfile
├── internal/
│   ├── database/                  # sqlc-generated DB layer
│   ├── handlers/                  # HTTP handlers
│   │   ├── category.go
│   │   ├── subscription.go
│   │   ├── user.go                # Auth handlers (register, login, etc.)
│   │   ├── user_categories.go
│   │   ├── user_profile.go
│   │   └── user_subscriptions.go
│   ├── middlewares/               # Logger, CORS, auth middleware
│   ├── models/                    # Request/response types, AppConfig
│   └── utils/                     # JWT, password hashing, JSON helpers
├── sql/
│   ├── queries/                   # sqlc input SQL queries
│   └── schema/                    # goose migration files
└── ui/                            # Svelte frontend
    └── src/
        ├── App.svelte             # Root component (layout, state, API)
        ├── app.css                # Global styles
        └── lib/
            ├── types.ts           # Shared TypeScript interfaces
            ├── helpers.ts         # Pure utility functions
            ├── AuthPage.svelte    # Login / register screen
            ├── Sidebar.svelte     # Navigation sidebar
            ├── DashboardPage.svelte
            ├── SubscriptionsPage.svelte
            ├── CategoriesPage.svelte
            ├── ProfilePage.svelte
            ├── SubModal.svelte    # Subscription add/edit modal
            └── CatModal.svelte    # Category add/edit modal
```

## License

MIT
