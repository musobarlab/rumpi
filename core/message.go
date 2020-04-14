package core

import (
	"time"
)

// MessageType type of message
type MessageType string

const (
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
	From        string      `json:"from,omitempty"`
	To          string      `json:"to,omitempty"`
	Content     string      `json:"content,omitempty"`
	Date        time.Time   `json:"date,omitempty"`
}
