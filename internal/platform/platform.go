// Package platform provides message platform interfaces
package platform

type MessagePlatformCapabilities struct {
	CanEditMessage    bool
	CanReplyToMessage bool
	CanSendMessage    bool
	CanThreadMessages bool
	MaxMessageLength  int
	MinMessageLength  int
}

type MessagePlatform interface {
	Start(stop chan bool)
	Capabilities() *MessagePlatformCapabilities
}
