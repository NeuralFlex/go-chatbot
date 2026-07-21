package services

import (
	"context"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
	"github.com/openai/openai-go/v3/responses"

	"go-chatbot/internal/models"
)

const systemPrompt = "You are a helpful assistant."

type OpenAIService struct {
	client openai.Client
	model  string
}

func OpenAIServiceInit(apiKey, model string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
		model:  model,
	}
}

func (s *OpenAIService) NewConversation(ctx context.Context) (string, error) {
	conv, err := s.client.Conversations.New(ctx, conversations.ConversationNewParams{})
	if err != nil {
		return "", err
	}
	return conv.ID, nil
}

func (s *OpenAIService) StreamReply(ctx context.Context, convID, query string) *ssestream.Stream[responses.ResponseStreamEventUnion] {
	params := responses.ResponseNewParams{
		Model: openai.ChatModel(s.model),
		Input: responses.ResponseNewParamsInputUnion{
			OfString: openai.String(query),
		},
		Conversation: responses.ResponseNewParamsConversationUnion{
			OfConversationObject: &responses.ResponseConversationParam{ID: convID},
		},
		Instructions: openai.String(systemPrompt),
	}
	return s.client.Responses.NewStreaming(ctx, params)
}

func (s *OpenAIService) History(ctx context.Context, convID string) ([]models.Message, error) {
	items, err := s.client.Conversations.Items.List(ctx, convID, conversations.ItemListParams{})
	if err != nil {
		return nil, err
	}

	var out []models.Message
	for _, item := range items.Data {
		if item.Type != "message" {
			continue
		}
		var sb strings.Builder
		for _, part := range item.Content.OfMessageContentArray {
			sb.WriteString(part.Text)
		}
		out = append(out, models.Message{Role: item.Role, Content: sb.String()})
	}
	return out, nil
}
