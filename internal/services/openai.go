package services

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"

	"go-chatbot/internal/models"
)

const systemPrompt = "You are a helpful assistant."

type OpenAIService struct {
	client openai.Client
	model  string
}

func OpenAIServiceInit(apiKey, model, baseURL string) *OpenAIService {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	return &OpenAIService{
		client: openai.NewClient(opts...),
		model:  model,
	}
}

func (s *OpenAIService) StreamReply(ctx context.Context, history []models.Message, query string) *ssestream.Stream[openai.ChatCompletionChunk] {
	msgs := make([]openai.ChatCompletionMessageParamUnion, 0, len(history)+2)
	msgs = append(msgs, openai.SystemMessage(systemPrompt))
	for _, m := range history {
		switch m.Role {
		case "user":
			msgs = append(msgs, openai.UserMessage(m.Content))
		case "assistant":
			msgs = append(msgs, openai.AssistantMessage(m.Content))
		}
	}
	msgs = append(msgs, openai.UserMessage(query))

	return s.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(s.model),
		Messages: msgs,
	})
}
