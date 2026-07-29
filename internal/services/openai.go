package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/conversations"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
	"github.com/openai/openai-go/v3/responses"

	"go-chatbot/internal/models"
)

const systemPrompt = "You are a helpful assistant."

const attachmentPrefix = "Attached file: "

func BuildMessageWithFile(filename, csvContent, userPrompt string) string {
	var b strings.Builder
	b.WriteString(attachmentPrefix)
	fmt.Fprintf(&b, "%d\n", len(filename))
	b.WriteString(filename)
	fmt.Fprintf(&b, "%d\n", len(csvContent))
	b.WriteString(csvContent)
	b.WriteString(userPrompt)
	return b.String()
}

func splitAttachment(content string) (filename, prompt string, ok bool) {
	rest, found := strings.CutPrefix(content, attachmentPrefix)
	if !found {
		return "", content, false
	}

	nameLen, rest, ok := readLengthPrefixed(rest)
	if !ok || nameLen > len(rest) {
		return "", content, false
	}
	filename, rest = rest[:nameLen], rest[nameLen:]

	csvLen, rest, ok := readLengthPrefixed(rest)
	if !ok || csvLen > len(rest) {
		return "", content, false
	}
	prompt = rest[csvLen:]
	return filename, prompt, true
}

func readLengthPrefixed(s string) (n int, rest string, ok bool) {
	i := strings.IndexByte(s, '\n')
	if i == -1 {
		return 0, "", false
	}
	n, err := strconv.Atoi(s[:i])
	if err != nil || n < 0 {
		return 0, "", false
	}
	return n, s[i+1:], true
}

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
	items, err := s.client.Conversations.Items.List(ctx, convID, conversations.ItemListParams{
		Order: conversations.ItemListParamsOrderAsc,
	})
	if err != nil {
		return nil, err
	}

	var out []models.Message
	for _, item := range items.Data {
		if item.Type != "message" {
			continue
		}
		parts := make([]string, len(item.Content.OfMessageContentArray))
		for i, part := range item.Content.OfMessageContentArray {
			parts[i] = part.Text
		}
		content := strings.Join(parts, " ")

		msg := models.Message{Role: item.Role, Content: content}
		if filename, prompt, ok := splitAttachment(content); ok {
			msg.Filename = filename
			msg.Content = prompt
		}
		out = append(out, msg)
	}
	return out, nil
}
