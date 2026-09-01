// Package openrouter provides an OpenRouter harness
package openrouter

import (
	"context"
	"flag"
	"log"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/models/operations"

	"github.com/haliphax/gobot/internal/harness"
)

var (
	OpenRouterAPIKey = flag.String("apiKey", "", "OpenRouter API key")
	Model            = flag.String("model", "xiaomi/mimo-v2.5", "Model identifier")
)

// ChatSender abstracts the OpenRouter chat Send method for testing
type ChatSender interface {
	Send(ctx context.Context,
		chatRequest components.ChatRequest,
		xOpenRouterMetadata *components.MetadataLevel,
		opts ...operations.Option,
	) (*operations.SendChatCompletionRequestResponse, error)
}

type OpenRouterClient struct {
	Chat  ChatSender
	Model *string
}

// New provides a new OpenRouterClient instance by reference
func New() *OpenRouterClient {
	s := openrouter.New(
		openrouter.WithSecurity(*OpenRouterAPIKey),
	)
	return &OpenRouterClient{s.Chat, Model}
}

// ProcessUserMessage sends the message to OpenRouter for generation
func (o *OpenRouterClient) ProcessUserMessage(message string) (string, error) {
	messageLen := len(message)
	snip := string(message[:min(20, messageLen)])
	if messageLen > 20 {
		snip += "..."
	}
	log.Printf("🤖 OpenRouter processing user message: %v", snip)

	ctx := context.Background()
	res, err := o.Chat.Send(ctx, components.ChatRequest{
		Model: o.Model,
		Messages: []components.ChatMessages{
			components.CreateChatMessagesUser(
				components.ChatUserMessage{
					Content: components.CreateChatUserMessageContentStr(message),
					Role:    components.ChatUserMessageRoleUser,
				},
			),
		},
	}, nil)
	if err != nil {
		log.Printf("❌ ERROR: %v", err.Error())
		return "", err
	}

	if res != nil {
		val, _ := res.ChatResult.Choices[0].Message.Content.GetOrZero()
		valStr := *val.Str
		valStrLen := len(valStr)
		snip = string(valStr[:min(20, valStrLen)])
		if valStrLen > 20 {
			snip += "..."
		}
		log.Printf("🗨️ OpenRouter response: %v", snip)
		return valStr, nil
	}

	return "", &harness.EmptyResponseError{}
}
