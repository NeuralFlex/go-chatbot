package models

import "time"

type ChatRequest struct {
	ConversationID string `form:"conversation_id"`
	Query          string `form:"query" binding:"required"`
}

type ConversationSummary struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	Role     string `json:"role"`
	Content  string `json:"content"`
	Filename string `json:"filename,omitempty"`
}

type ConversationDetail struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Messages []Message `json:"messages"`
}

type Preset struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Text  string `json:"text"`
}
