package main

import (
  "log"

  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"

  _ "go-chatbot/docs"
  "go-chatbot/internal/config"
  "go-chatbot/internal/database"
  "go-chatbot/internal/handlers"
  "go-chatbot/internal/repository"
  "go-chatbot/internal/routes"
  "go-chatbot/internal/services"
)

// @title        Go Chatbot API
// @version      1.0
// @description  API documentation for the go-chatbot service.
func main() {
  _ = godotenv.Load() // ignore error: .env may be absent (e.g. Docker sets env vars directly)

  cfg, err := config.Load()
  if err != nil {
    log.Fatalf("config: %v", err)
  }

  db, err := database.Open(cfg.DBPath)
  if err != nil {
    log.Fatalf("database: %v", err)
  }
  defer db.Close()

  repo := &repository.ConversationRepository{DB: db.GormDB}
  openAI := services.OpenAIServiceInit(cfg.OpenAIKey, cfg.Model, cfg.BaseURL)
  conv := &handlers.ConversationsHandler{Repo: repo, OpenAI: openAI}

  router := gin.Default()
  routes.Register(router, conv)
  router.Run(":" + cfg.Port)
}