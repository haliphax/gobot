// Package harness - OpenRouter client
package harness

import (
	"context"
	"flag"
	"log"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
)

var (
	OpenRouterAPIKey = flag.String("apiKey", "", "OpenRouter API key")
	Model            = "xiaomi/mimo-v2.5"
)

type OpenRouterClient struct{}

// ProcessUserMessage - Generate response to user message
func (o *OpenRouterClient) ProcessUserMessage(message string) (string, error) {
	ctx := context.Background()

	s := openrouter.New(
		openrouter.WithSecurity(*OpenRouterAPIKey),
	)

	res, err := s.Chat.Send(ctx, components.ChatRequest{
		Model: &Model,
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
		log.Fatal(err)
	}

	if res != nil {
		val, _ := res.ChatResult.Choices[0].Message.Content.GetOrZero()
		return *val.Str, nil
	}

	return "", &EmptyResponseError{}
}
