package chathub

import (
	"time"
)

// OnlineUser represent online user
type OnlineUser struct {
	Username string `json:"username,omitempty"`
	Status   bool   `json:"status"`
}

// Message structure
type Message struct {
	MessageType MessageType   `json:"messageType,omitempty"`
	AuthKey     string        `json:"authKey,omitempty"`
	Token       string        `json:"token,omitempty"`
	Username    string        `json:"username,omitempty"`
	From        string        `json:"from,omitempty"`
	To          string        `json:"to,omitempty"`
	Content     string        `json:"content,omitempty"`
	Date        time.Time     `json:"date,omitempty"`
	OnlineUsers []*OnlineUser `json:"onlineUsers,omitempty"`
}
