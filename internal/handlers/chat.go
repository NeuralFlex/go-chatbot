package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-chatbot/internal/middleware"
	"go-chatbot/internal/models"
	"go-chatbot/internal/utils"
)

const maxFileSize = 300 << 10 // 300KB — keeps the inlined CSV well under the model's context window

// Send godoc
// @Summary      Send a message (streams the reply as SSE)
// @Description  Omit conversation_id to start a new chat; include it to continue one. A .csv file (optional, 300KB max) may only be attached when starting a new chat.
// @Tags         chat
// @Accept       multipart/form-data
// @Produce      text/event-stream
// @Param        X-User-ID        header  string  true   "User ID"
// @Param        query            formData  string  true  "Message"
// @Param        conversation_id  formData  string  false "Conversation ID"
// @Param        file             formData  file    false "CSV attachment (300KB max)"
// @Success      200        {string}  string  "SSE stream"
// @Failure      400        {object}  map[string]string
// @Failure      404        {object}  map[string]string
// @Router       /chat [post]
func (h *ConversationsHandler) Send(c *gin.Context) {
	userID := middleware.UserID(c)
	ctx := c.Request.Context()

	var req models.ChatRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fileContent string
	var filename string
	if header, err := c.FormFile("file"); err == nil {
		if req.ConversationID != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "a file can only be attached when starting a new conversation"})
			return
		}
		if !strings.HasSuffix(strings.ToLower(header.Filename), ".csv") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only .csv files are supported"})
			return
		}
		if header.Size > maxFileSize {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file exceeds 300KB limit"})
			return
		}
		f, err := header.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		filename = header.Filename
		fileContent = string(data)
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
	} else {
		found, _, _, err := h.Repo.Owns(ctx, userID, convID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
			return
		}
	}

	message := req.Query
	if filename != "" {
		message = utils.BuildMessageWithFile(filename, fileContent, req.Query)
	}

	stream := h.OpenAI.StreamReply(ctx, convID, message)
	defer stream.Close()

	firstEvent := isNew
	c.Stream(func(w io.Writer) bool {
		if !stream.Next() {
			if err := stream.Err(); err != nil {
				c.SSEvent("error", err.Error())
			} else {
				c.SSEvent("done", "[DONE]")
			}
			return false
		}
		if firstEvent {
			firstEvent = false
			if err := h.Repo.Create(ctx, userID, convID, truncateTitle(req.Query), filename); err != nil {
				c.SSEvent("error", err.Error())
				return false
			}
			c.SSEvent("conversation", convID)
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
	return query[:maxLen] + "..."
}
