package models

type ChatRequest struct {
	ConversationID string `json:"conversation_id,omitempty"`
	Query          string `json:"query" binding:"required"`
}

type ConversationSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ConversationDetail struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Messages []Message `json:"messages"`
}
