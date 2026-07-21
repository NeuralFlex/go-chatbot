package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

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

	convID := req.ConversationID
	isNew := convID == ""

	if isNew {
		var err error
		convID, err = h.OpenAI.NewConversation(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := h.Repo.Create(ctx, userID, convID, truncateTitle(req.Query)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		found, _, err := h.Repo.Owns(ctx, userID, convID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}
	}

	stream := h.OpenAI.StreamReply(ctx, convID, req.Query)
	defer stream.Close()

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
			} else {
				c.SSEvent("done", "[DONE]")
			}
			return false
		}
		event := stream.Current()
		if event.Type == "response.output_text.delta" {
			c.SSEvent("message", event.Delta)
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
