// Package openrouter provides an OpenRouter harness
package openrouter

import (
	"context"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/models/operations"

	"github.com/haliphax/gobot/internal/harness"
)

type OpenRouterConfig struct {
	Token string
}

type Config struct {
	OpenRouter OpenRouterConfig
}

// ChatSender abstracts the OpenRouter chat Send method for testing
type ChatSender interface {
	Send(ctx context.Context,
		chatRequest components.ChatRequest,
		xOpenRouterMetadata *components.MetadataLevel,
		opts ...operations.Option,
	) (*operations.SendChatCompletionRequestResponse, error)
}

type OpenRouterClient struct {
	OpenRouterConfig
	Chat         ChatSender
	currentModel string
}

// New provides a new OpenRouterClient instance by reference
func New(configFilename *string) *OpenRouterClient {
	file, err := os.ReadFile(*configFilename)
	if err != nil {
		panic(err)
	}

	var baseConf Config
	_, err = toml.Decode(string(file), &baseConf)
	if err != nil {
		panic(err)
	}

	conf := baseConf.OpenRouter
	s := openrouter.New(
		openrouter.WithSecurity(conf.Token),
	)
	return &OpenRouterClient{conf, s.Chat, ""}
}

func (o *OpenRouterClient) Model() string {
	return o.currentModel
}

func (o *OpenRouterClient) SetModel(model string) string {
	o.currentModel = model
	return o.Model()
}

// ProcessUserMessage sends the message to OpenRouter for generation
func (o *OpenRouterClient) ProcessUserMessage(message string) (string, error) {
	messageLen := len(message)
	snip := string(message[:min(harness.UserMessageLogSnippetLength, messageLen)])
	if messageLen > harness.UserMessageLogSnippetLength {
		snip += "..."
	}
	log.Printf("🤖 OpenRouter processing user message (model: %v): %v", o.Model(), snip)

	ctx := context.Background()
	res, err := o.Chat.Send(ctx, components.ChatRequest{
		Model: new(o.Model()),
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
		snip = string(valStr[:min(harness.AgentMessageLogSnippetLength, valStrLen)])
		if valStrLen > harness.AgentMessageLogSnippetLength {
			snip += "..."
		}
		log.Printf("🗨️ OpenRouter response (model: %v): %v", o.Model(), snip)
		return valStr, nil
	}

	return "", &harness.EmptyResponseError{}
}
