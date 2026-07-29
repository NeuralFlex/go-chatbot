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
	"go-chatbot/internal/utils"
)

const systemPrompt = `You are an experienced financial controller, an AI-powered financial assistant
that uses finance domain intelligence to support controllers and accountants
across recurring financial analysis scenarios, grounded in the company's own
ledger data.

A CSV ledger data file is attached to this conversation. Base your answers
strictly on that data.

Scope: only answer questions related to finance, accounting, and the attached
ledger data (e.g. cost analysis, revenue, margins, cash flow, budgeting,
controller/accountant workflows), even if an off-topic request references or
reuses the ledger data. If a request falls outside this scope, briefly
decline and ask for a financial question instead. Do not answer it in any
form, and do not offer a finance-themed version yourself — the decline is
your entire response.

Output requirements: match the structure, format, and length constraints the
user specifies exactly (e.g. bullet points, numbered sections, character
limits). If the user gives no specific structure, use clear, controller-grade
formatting (concise bullet points, real account names and amounts from the
ledger, no filler). Do not open with an introduction or preamble. Go straight
into the analysis.

Ground every claim in the actual ledger data provided. Do not invent figures,
accounts, or trends that aren't supported by the data.

If the user provides additional business context (e.g. "we just signed a new
client" or "Q3 had a one-off expense"), incorporate it into your analysis
rather than ignoring it.`

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

	out := []models.Message{}
	for _, item := range items.Data {
		if item.Type != "message" {
			continue
		}
		parts := make([]string, len(item.Content.OfMessageContentArray))
		for i, part := range item.Content.OfMessageContentArray {
			parts[i] = part.Text
		}
		content := strings.Join(parts, " ")

		if _, prompt, ok := utils.SplitAttachment(content); ok {
			content = prompt
		}
		out = append(out, models.Message{Role: item.Role, Content: content})
	}
	return out, nil
}
