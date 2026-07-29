# go-chatbot

Backend API for an AI-powered financial controller assistant chatbot, built in Go and
Gin, using the OpenAI SDK.

## Setup

1. Install dependencies:

   ```bash
   go mod download
   ```

2. Start a local Postgres database:

   ```bash
   docker compose up -d
   ```

   Stop it with `docker compose down` (add `-v` to also delete the data).

3. Configure environment:

   ```bash
   cp .env.example .env
   ```

   Then fill in `OPENAI_API_KEY` in `.env`. The default `DATABASE_URL` already
   matches the `docker-compose.yml` credentials.

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
| `internal/database` | Postgres connection setup, goose migrations, GORM auto-migrate |
| `internal/models` | DB models and DTOs |
| `internal/repository` | Conversation data access |
| `internal/services` | OpenAI SDK integration, system prompt, preset prompts |
| `internal/utils` | Attachment message encoding/parsing |
| `internal/handlers` | HTTP request handlers |
| `internal/routes` | Route registration |
| `internal/middleware` | Gin middleware (user identification) |
| `docs` | Generated Swagger docs |
