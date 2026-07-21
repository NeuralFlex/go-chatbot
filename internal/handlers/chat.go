package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/ssestream"

	"go-chatbot/internal/middleware"
	"go-chatbot/internal/models"
)

// Send godoc
// @Summary      Send a message (streams the reply as SSE)
// @Description  Omit conversation_id to start a new chat; include it to continue one.
// @Tags         chat
// @Accept       json
// @Produce      text/event-stream
// @Param        X-User-ID  header  string              true  "User ID"
// @Param        request    body    models.ChatRequest  true  "Message"
// @Success      200        {string}  string  "SSE stream"
// @Failure      400        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Router       /chat [post]
func (h *ConversationsHandler) Send(c *gin.Context) {
	userID := middleware.UserID(c)
	ctx := c.Request.Context()

	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	convID, history, isNew, status, err := h.resolveConversation(ctx, userID, req)
	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	stream := h.OpenAI.StreamReply(ctx, history, req.Query)
	defer stream.Close()

	h.streamAndSave(c, stream, convID, req.Query, isNew)
}

// resolveConversation creates a new conversation or loads an existing one's history.
func (h *ConversationsHandler) resolveConversation(ctx context.Context, userID string, req models.ChatRequest) (convID string, history []models.Message, isNew bool, status int, err error) {
	if req.ConversationID == "" {
		convID = "conv_" + uuid.NewString()
		if err := h.Repo.Create(ctx, userID, convID, truncateTitle(req.Query)); err != nil {
			return "", nil, false, http.StatusInternalServerError, err
		}
		return convID, nil, true, 0, nil
	}

	convID = req.ConversationID
	found, _, err := h.Repo.Owns(ctx, userID, convID)
	if err != nil {
		return "", nil, false, http.StatusInternalServerError, err
	}
	if !found {
		return "", nil, false, http.StatusNotFound, errors.New("conversation not found")
	}

	history, err = h.Repo.Messages(ctx, convID)
	if err != nil {
		return "", nil, false, http.StatusInternalServerError, err
	}
	return convID, history, false, 0, nil
}

// streamAndSave relays the reply to the client as SSE and saves the turn once streaming ends.
func (h *ConversationsHandler) streamAndSave(c *gin.Context, stream *ssestream.Stream[openai.ChatCompletionChunk], convID, query string, isNew bool) {
	ctx := c.Request.Context()
	var reply strings.Builder
	firstEvent := isNew

	c.Stream(func(w io.Writer) bool {
		if firstEvent {
			firstEvent = false
			c.SSEvent("conversation", convID)
			return true
		}

		if !stream.Next() {
			if err := stream.Err(); err != nil {
				c.SSEvent("error", err.Error())
				return false
			}
			if err := h.Repo.AddMessages(ctx, convID,
				models.Message{Role: "user", Content: query},
				models.Message{Role: "assistant", Content: reply.String()},
			); err != nil {
				c.SSEvent("error", err.Error())
				return false
			}
			c.SSEvent("done", "[DONE]")
			return false
		}

		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			delta := chunk.Choices[0].Delta.Content
			reply.WriteString(delta)
			c.SSEvent("message", delta)
		}
		return true
	})
}

func truncateTitle(query string) string {
	const maxLen = 50
	if len(query) <= maxLen {
		return query
	}
	return query[:maxLen]
}
