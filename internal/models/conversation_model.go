package models

type Conversation struct {
	ConversationID string `gorm:"primaryKey"`
	UserID         string `gorm:"index;not null"`
	Title          string `gorm:"not null"`
	CreatedAt      string `gorm:"not null"`
}

type ConversationMessage struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	ConversationID string `gorm:"index;not null"`
	Role           string `gorm:"not null"`
	Content        string `gorm:"not null"`
	CreatedAt      string `gorm:"not null"`
}

