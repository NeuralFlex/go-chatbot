# go-chatbot

Backend API chatbot built in Go and Gin, using the OpenAI SDK.

## Setup

1. Install dependencies:

   ```bash
   go mod download
   ```

2. Configure environment:

   ```bash
   cp .env.example .env
   ```

   Then fill in `OPENAI_API_KEY` in `.env`.

## Run

```bash
go run ./cmd/server
```

Server starts on `PORT` (default `8080`). Swagger docs are available at `/swagger/index.html`.

## Folder structure

| Path | Description |
|---|---|
| `cmd/server` | Entrypoint (`main.go`) |
| `internal/config` | Env/config loading |
| `internal/database` | SQLite connection setup |
| `internal/models` | DB models and DTOs |
| `internal/repository` | Conversation data access |
| `internal/services` | OpenAI SDK integration |
| `internal/handlers` | HTTP request handlers |
| `internal/routes` | Route registration |
| `internal/middleware` | Gin middleware (user identification) |
| `docs` | Generated Swagger docs |
