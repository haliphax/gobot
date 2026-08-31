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
	Model            = flag.String("model", "xiaomi/mimo-v2.5", "Model identifier")
)

type OpenRouterClient struct {
	Session *openrouter.OpenRouter
	Model   string
}

func (o *OpenRouterClient) Init() *OpenRouterClient {
	s := openrouter.New(
		openrouter.WithSecurity(*OpenRouterAPIKey),
	)
	o.Session = s
	return o
}

// ProcessUserMessage - Generate response to user message
func (o *OpenRouterClient) ProcessUserMessage(message string) (string, error) {
	ctx := context.Background()

	res, err := o.Session.Chat.Send(ctx, components.ChatRequest{
		Model: Model,
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
		return *val.Str, nil
	}

	return "", &EmptyResponseError{}
}
