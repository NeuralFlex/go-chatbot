package models

import "time"

type Conversation struct {
	ConversationID string `gorm:"primaryKey"`
	UserID         string `gorm:"index;not null"`
	Title          string `gorm:"not null"`
	Filename       string
	CreatedAt      time.Time
}
