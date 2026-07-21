package models

type Conversation struct {
	ConversationID string `gorm:"primaryKey"`
	UserID         string `gorm:"index;not null"`
	Title          string `gorm:"not null"`
	CreatedAt      string `gorm:"not null"`
}
