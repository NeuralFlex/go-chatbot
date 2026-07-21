package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"go-chatbot/internal/models"
)

type ConversationRepository struct {
	DB *gorm.DB
}

func (r *ConversationRepository) Create(ctx context.Context, userID, convID, title string) error {
	row := models.Conversation{
		ConversationID: convID,
		UserID:         userID,
		Title:          title,
		CreatedAt:      time.Now().UTC().Format(time.RFC3339),
	}
	return r.DB.WithContext(ctx).Create(&row).Error
}

func (r *ConversationRepository) List(ctx context.Context, userID string) ([]models.ConversationSummary, error) {
	var rows []models.Conversation
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	summaries := make([]models.ConversationSummary, len(rows))
	for i, row := range rows {
		summaries[i] = models.ConversationSummary{
			ID:        row.ConversationID,
			Title:     row.Title,
			CreatedAt: row.CreatedAt,
		}
	}
	return summaries, nil
}


func (r *ConversationRepository) Owns(ctx context.Context, userID, convID string) (found bool, title string, err error) {
	var row models.Conversation
	dbErr := r.DB.WithContext(ctx).
		First(&row, "conversation_id = ? AND user_id = ?", convID, userID).Error
	if errors.Is(dbErr, gorm.ErrRecordNotFound) {
		return false, "", nil
	}
	if dbErr != nil {
		return false, "", dbErr
	}
	return true, row.Title, nil
}

func (r *ConversationRepository) Messages(ctx context.Context, convID string) ([]models.Message, error) {
	var rows []models.ConversationMessage
	err := r.DB.WithContext(ctx).
		Where("conversation_id = ?", convID).
		Order("id ASC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	messages := make([]models.Message, len(rows))
	for i, row := range rows {
		messages[i] = models.Message{Role: row.Role, Content: row.Content}
	}
	return messages, nil
}

func (r *ConversationRepository) AddMessages(ctx context.Context, convID string, msgs ...models.Message) error {
	now := time.Now().UTC().Format(time.RFC3339)
	rows := make([]models.ConversationMessage, len(msgs))
	for i, m := range msgs {
		rows[i] = models.ConversationMessage{
			ConversationID: convID,
			Role:           m.Role,
			Content:        m.Content,
			CreatedAt:      now,
		}
	}
	return r.DB.WithContext(ctx).Create(&rows).Error
}
