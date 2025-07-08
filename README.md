# Chirpy

Chirpy is a Go-based web service for managing users, chirps (short messages), and authentication. It features a modular structure, RESTful endpoints, and concurrency-safe metrics. Basically a twitter backend clone.

## Features
- User registration and authentication (API keys, JWT, password)
- Chirp creation, retrieval, and deletion
- Admin endpoints for metrics and readiness
- Static file serving
- SQL migrations and queries
- Concurrency-safe hit counter using atomic operations

## Project Structure
```
chirpy/
  cmd/                # Main application entry point
  internal/
    app/              # Application and environment setup
    auth/             # Authentication logic
    config/           # Config and middleware
    database/         # Database models and queries
    handlers/         # HTTP handlers (admin, chirps, hooks, tokens, users)
    ...
  sql/
    queries/          # SQL query files
    schema/           # Database schema migrations
  go.mod, go.sum      # Go module files
  sqlc.yaml           # SQLC configuration
```

## Getting Started

### Prerequisites
- Go 1.19+
- PostgreSQL

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/goinginblind/chirpy.git
   cd chirpy
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up the database:
   - Create a PostgreSQL database.
   - Run migrations in `sql/schema/` (e.g., using [goose](https://github.com/pressly/goose)):
     ```bash
     goose -dir sql/schema postgres "<connection-string>" up
     ```
   - Enable the `uuid-ossp` extension:
     ```sql
     CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
     ```
4. Build and run the server:
   ```bash
   go run ./cmd/chirpy/main.go
   ```

### Environment Variables
- Configure environment variables as needed (e.g., database connection string, port).

## Usage
- Access the app at `http://localhost:8080/app/`

## API Endpoints

Below are the main API endpoints, their methods, and example request/response formats:

### Health & Admin

- **GET /api/healthz**
  - Description: Health check endpoint.
  - Response: `200 OK`, plain text "OK"

- **GET /admin/metrics**
  - Description: Returns server metrics (e.g., file server hits).
  - Response: `200 OK`, plain text (e.g., `Hits: 42`)

- **POST /admin/reset**
  - Description: Resets users and file server hit counter (dev only).
  - Response: `200 OK`, JSON: `{ "message": "users reset, amount of hits set to 0" }`

### Users

- **POST /api/users**
  - Description: Register a new user.
  - Request JSON: `{ "email": "user@example.com", "password": "secret" }`
  - Response: `201 Created`, JSON: `{ "id": "...", "email": "..." }`

- **POST /api/login**
  - Description: User login.
  - Request JSON: `{ "email": "user@example.com", "password": "secret" }`
  - Response: `200 OK`, JSON: `{ "token": "..." }`

- **PUT /api/users**
  - Description: Change user login info (e.g., password).
  - Request header should include an access token
  - Request JSON: `{ "password": "newpassword" }`
  - Response: `200 OK`, JSON: `{ "message": "updated" }`

### Tokens

- **POST /api/refresh**
  - Description: Refresh access token.
  - Request header should include a refresh token
  - Response: `200 OK`, JSON: `{ "token": "..." }`

- **POST /api/revoke**
  - Description: Revoke a refresh token.
  - Request header should include a refresh token
  - Response: `200 OK`, JSON: `{ "message": "revoked" }`

### Chirps

- **POST /api/chirps**
  - Description: Create a new chirp.
  - Request JSON: `{ "body": "Hello world!" }`
  - Request header should include an access token
  - Response: `201 Created`, JSON: `{ "id": "...", "body": "...", "created_at": "...", "updated_at": "...", "user_id": "..."}`

- **GET /api/chirps**
  - Description: Get all chirps; author ID and sort by publish date are optional and can be provided as URL queries.
  - Response: `200 OK`, JSON: `[ { "id": "...", "body": "...", "created_at": "...", "updated_at": "...", "user_id": "..."}, ... ]`

- **GET /api/chirps/{chirpID}**
  - Description: Get a single chirp by ID.
  - Response: `200 OK`, JSON: `{ "id": "...", "body": "...", "created_at": "...", "updated_at": "...", "user_id": "..."}`

- **DELETE /api/chirps/{chirpID}**
  - Description: Delete a chirp by ID.
  - Request header should include an access token
  - Response: `204 No Content`

### Webhooks

- **POST /api/polka/webhooks**
  - Description: Upgrade user to Chirpy Red (webhook endpoint).
  - Request header should have a valid API key in it
  - Response: `204 No Content`

---

For more details, see the handler implementations in `internal/handlers/`.

## Testing
Run tests with (altough there are only the auth package unit tests):
```bash
go test ./...
```

## License
MIT
