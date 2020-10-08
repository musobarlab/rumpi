package chathub

import "time"

// MessageType type of message
type MessageType string

const (
	// AuthMessage message
	AuthMessage MessageType = "authMessage"

	// AuthFail message
	AuthFail MessageType = "authFail"

	// UserJoined message
	UserJoined MessageType = "userJoined"

	// PrivateMessage message
	PrivateMessage MessageType = "privateMessage"

	//Broadcast message
	Broadcast MessageType = "broadcast"

	//UsersStatus message
	UsersStatus MessageType = "usersStatus"

	// WriteWait Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// PongWait Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// PongPeriod should less than PongWait
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

// Clients type alias of map Client
type Clients map[string]*Client
