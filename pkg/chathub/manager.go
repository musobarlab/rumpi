package chathub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/musobarlab/rumpi/pkg/jwt"
)

const (
	// UpdateStatusWait to update users status
	UpdateStatusWait = 2 * time.Second
)

// Manager will manage client and chat
type Manager struct {
	AuthKey         string
	Clients         map[string]*Client
	Register        chan *Client
	Unregister      chan *Client
	AuthSuccess     chan *Client
	IncomingMessage chan *Message
	Upgrader        websocket.Upgrader
	JwtService      jwt.JwtService
	sync.RWMutex
}

// NewManager function
func NewManager(authKey string, jwtService jwt.JwtService) *Manager {
	clients := make(map[string]*Client)
	incomingMessage := make(chan *Message)
	register := make(chan *Client)
	unregister := make(chan *Client)
	authSuccess := make(chan *Client)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}

	return &Manager{
		AuthKey:         authKey,
		Clients:         clients,
		IncomingMessage: incomingMessage,
		Register:        register,
		Unregister:      unregister,
		AuthSuccess:     authSuccess,
		Upgrader:        upgrader,
		JwtService:      jwtService,
	}

}

// Handle function will handle incoming client
func (manager *Manager) Handle() {
	ticker := time.NewTicker(UpdateStatusWait)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case client := <-manager.AuthSuccess:

			// create room for client using its username
			// client.joinRoom(client.Username)

			manager.addClient(client.Username, client)

			message := &Message{
				From:        client.Username,
				MessageType: Broadcast,
				Date:        time.Now(),
				Content:     "has join the chat",
			}
			manager.send(message, client)

		case <-manager.Register:

		case client := <-manager.Unregister:

			// set client online status to false
			client.IsOnline = false

			message := &Message{
				From:        client.Username,
				MessageType: Broadcast,
				Date:        time.Now(),
				Content:     "has leave the chat",
			}
			manager.send(message, client)
			manager.deleteClient(client.Username)

		case m := <-manager.IncomingMessage:
			for client := range manager.Clients {
				fmt.Println(client)
			}

			switch m.MessageType {
			case PrivateMessage:
				// send to specific user
				manager.sendPrivate(m)

			case Broadcast:
				// send to every client that is currently connected
				manager.send(m, nil)
			}
		case <-ticker.C:
			var users []*OnlineUser
			for _, client := range manager.Clients {
				users = append(users, &OnlineUser{Username: client.Username, Status: client.IsOnline})
			}
			msg := Message{
				MessageType: UsersStatus,
				OnlineUsers: users,
			}

			manager.send(&msg, nil)
		}
	}
}

func (manager *Manager) send(message *Message, ignore *Client) {
	msg, _ := json.Marshal(message)

	manager.RLock()
	defer manager.RUnlock()

	for _, client := range manager.Clients {
		if client != ignore {
			select {
			case client.MsgChan <- msg:
			default:
				close(client.MsgChan)
				manager.deleteClient(client.Username)
			}
		}
	}
}

func (manager *Manager) sendPrivate(message *Message) {
	msg, _ := json.Marshal(message)

	manager.RLock()
	defer manager.RUnlock()

	if client, ok := manager.Clients[message.To]; ok {
		select {
		case client.MsgChan <- msg:
		default:
			close(client.MsgChan)
			manager.deleteClient(client.Username)
		}
	}
}

//addClient function will push new client to the map clients
func (p *Manager) addClient(key string, client *Client) {
	p.Lock()
	p.Clients[key] = client
	p.Unlock()
}

//deleteClient function will delete client by specific key from map clients
func (p *Manager) deleteClient(key string) {
	p.Lock()
	delete(p.Clients, key)
	p.Unlock()
}
