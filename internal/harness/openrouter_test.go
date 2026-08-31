// Package harness - OpenRouter tests
package harness

import (
	"context"
	"testing"

	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/OpenRouterTeam/go-sdk/models/operations"
	"github.com/OpenRouterTeam/go-sdk/optionalnullable"
)

type MockChatSender struct {
	SendFn func(ctx context.Context,
		chatRequest components.ChatRequest,
		xOpenRouterMetadata *components.MetadataLevel,
		opts ...operations.Option,
	) (*operations.SendChatCompletionRequestResponse, error)
}

func (m *MockChatSender) Send(ctx context.Context,
	chatRequest components.ChatRequest,
	xOpenRouterMetadata *components.MetadataLevel,
	opts ...operations.Option,
) (*operations.SendChatCompletionRequestResponse, error) {
	return m.SendFn(ctx, chatRequest, xOpenRouterMetadata, opts...)
}

func TestProcessUserMessageCallsSend(t *testing.T) {
	Model = new("test-model")
	var capturedRequest components.ChatRequest
	sendCalled := false

	mock := &MockChatSender{
		SendFn: func(ctx context.Context,
			chatRequest components.ChatRequest,
			xOpenRouterMetadata *components.MetadataLevel,
			opts ...operations.Option,
		) (*operations.SendChatCompletionRequestResponse, error) {
			sendCalled = true
			capturedRequest = chatRequest

			resp := operations.CreateSendChatCompletionRequestResponseChatResult(
				components.ChatResult{
					Choices: []components.ChatChoice{
						{
							Message: components.ChatAssistantMessage{
								Content: optionalnullable.From(
									&components.ChatAssistantMessageContent{
										Str: new("I'm a mock response!"),
									},
								),
							},
						},
					},
				},
			)

			return &resp, nil
		},
	}

	client := &OpenRouterClient{
		Chat:  mock,
		Model: new("test-model"),
	}

	result, err := client.ProcessUserMessage("hello")
	if err != nil {
		t.Fatal(err)
	}

	if !sendCalled {
		t.Error("Send was never called")
	}

	if result != "I'm a mock response!" {
		t.Errorf("got %q, want %q", result, "I'm a mock response!")
	}

	// Verify the message was passed through
	if len(capturedRequest.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(capturedRequest.Messages))
	}

	if capturedRequest.Model != nil && *capturedRequest.Model != "test-model" {
		t.Errorf("expected model 'test-model', got %q", *capturedRequest.Model)
	}
}
