package core

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Manager model
type Manager struct {
	Clients         map[*Client]bool
	Register        chan *Client
	Unregister      chan *Client
	IncomingMessage chan *Message
	Upgrader        websocket.Upgrader
	sync.RWMutex
}

// NewManager function
func NewManager() *Manager {
	clients := make(map[*Client]bool)
	incomingMessage := make(chan *Message)
	register := make(chan *Client)
	unregister := make(chan *Client)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}

	return &Manager{
		Clients:         clients,
		IncomingMessage: incomingMessage,
		Register:        register,
		Unregister:      unregister,
		Upgrader:        upgrader,
	}

}

// Run function will run Manager process
func (manager *Manager) Run() {
	for {
		select {
		case client := <-manager.Register:
			message := &Message{
				From:        client.ID,
				MessageType: Broadcast,
				Date:        time.Now(),
				Content:     "has join the chat",
			}
			manager.send(message, client)
			manager.AddClient(client, true)
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client]; ok {
				message := &Message{
					From:        client.ID,
					MessageType: Broadcast,
					Date:        time.Now(),
					Content:     "has leave the chat",
				}
				manager.send(message, client)
				manager.DeleteClient(client)
			}
		case m := <-manager.IncomingMessage:

			switch m.MessageType {
			case PrivateMessage:
				// send to specific user
				manager.sendPrivate(m)
			case Broadcast:
				// send to every client that is currently connected
				manager.send(m, nil)
			}
		}
	}
}

func (manager *Manager) send(message *Message, ignore *Client) {
	msg, _ := json.Marshal(message)
	for client := range manager.Clients {
		if client != ignore {
			select {
			case client.MsgChan <- msg:
			default:
				close(client.MsgChan)
				manager.DeleteClient(client)
			}
		}
	}
}

func (manager *Manager) sendPrivate(message *Message) {
	msg, _ := json.Marshal(message)
	for client := range manager.Clients {
		if _, ok := client.Room[message.To]; ok {
			select {
			case client.MsgChan <- msg:
			default:
				close(client.MsgChan)
				manager.DeleteClient(client)
			}
		}
	}
}

//AddClient function will push new client to the map clients
func (p *Manager) AddClient(key *Client, b bool) {

	// add default room with its ID
	key.AddRoom(key.ID)

	p.Lock()
	p.Clients[key] = b
	p.Unlock()
}

//DeleteClient function will delete client by specific key from map clients
func (p *Manager) DeleteClient(key *Client) {
	close(key.MsgChan)
	p.Lock()
	delete(p.Clients, key)
	p.Unlock()
}
