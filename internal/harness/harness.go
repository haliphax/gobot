// Package harness
package harness

type EmptyResponseError struct{}

func (e *EmptyResponseError) Error() string {
	return "Empty response"
}

type ModelProviderClient interface {
	ProcessUserMessage(message string) (string, error)
}
