package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-chatbot/internal/middleware"
	"go-chatbot/internal/models"
	"go-chatbot/internal/repository"
	"go-chatbot/internal/services"
)

type ConversationsHandler struct {
	Repo   *repository.ConversationRepository
	OpenAI *services.OpenAIService
}

// List godoc
// @Summary      List the user's conversations
// @Tags         conversations
// @Produce      json
// @Param        X-User-ID  header    string  true  "User ID"
// @Success      200        {array}   models.ConversationSummary
// @Router       /conversations [get]
func (h *ConversationsHandler) List(c *gin.Context) {
	userID := middleware.UserID(c)

	summaries, err := h.Repo.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summaries)
}

// Get godoc
// @Summary      Open a conversation and load its history
// @Tags         conversations
// @Produce      json
// @Param        id         path      string  true  "Conversation ID"
// @Param        X-User-ID  header    string  true  "User ID"
// @Success      200        {object}  models.ConversationDetail
// @Failure      404        {object}  map[string]string
// @Router       /conversations/{id} [get]
func (h *ConversationsHandler) Get(c *gin.Context) {
	userID := middleware.UserID(c)
	convID := c.Param("id")
	ctx := c.Request.Context()

	found, title, err := h.Repo.Owns(ctx, userID, convID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	messages, err := h.Repo.Messages(ctx, convID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.ConversationDetail{
		ID:       convID,
		Title:    title,
		Messages: messages,
	})
}
