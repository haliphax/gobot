// Package harness provides agent harness functionality
package harness

const (
	UserMessageLogSnippetLength  = 120
	AgentMessageLogSnippetLength = 120
)

type EmptyResponseError struct{}

func (e *EmptyResponseError) Error() string {
	return "Empty response"
}

type ModelProviderClient interface {
	ProcessUserMessage(message string) (string, error)
	// props
	Model() string
	SetModel(model string) string
}
