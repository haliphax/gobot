// Package platform provides message platform interfaces
package platform

type MessagePlatform interface {
	Start(stop chan bool)
}
