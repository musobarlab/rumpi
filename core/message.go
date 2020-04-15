package core

import (
	"time"
)

// MessageType type of message
type MessageType string

const (
	// AuthMessage message
	AuthMessage MessageType = "authMessage"

	// UserJoined message
	UserJoined MessageType = "userJoined"

	// PrivateMessage message
	PrivateMessage MessageType = "privateMessage"

	//Broadcast message
	Broadcast MessageType = "broadcast"
)

// Message structure
type Message struct {
	MessageType MessageType `json:"messageType,omitempty"`
	AuthKey     string      `json:"authKey,omitempty"`
	Username    string      `json:"username,omitempty"`
	From        string      `json:"from,omitempty"`
	To          string      `json:"to,omitempty"`
	Content     string      `json:"content,omitempty"`
	Date        time.Time   `json:"date,omitempty"`
}
