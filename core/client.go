package core

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientModel struct {
	Username string `json:"username"`
}

// Client model
type Client struct {
	ID      string
	Conn    *websocket.Conn
	MsgChan chan []byte
	Manager *Manager
	Room    map[string]bool
	sync.RWMutex
}

//AddRoom function will push new room to the map rooms
func (client *Client) AddRoom(key string) {
	client.Lock()
	client.Room[key] = true
	client.Unlock()
}

//DeleteRoom function will delete room by specific key from map rooms
func (client *Client) DeleteRoom(key string) {
	client.Lock()
	delete(client.Room, key)
	client.Unlock()
}

// Read function will read incoming message from client
func (c *Client) Read() {

	// if function finish, remove client
	// and also close connection
	defer func() {
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// read message from client
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			c.Manager.Unregister <- c
			c.Conn.Close()
			break
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			c.Manager.Unregister <- c
			c.Conn.Close()
			break
		}

		message.From = c.ID
		message.Date = time.Now()

		// send message to Manager's IncomingMessage
		c.Manager.IncomingMessage <- &message

	}

}

// Write function will write message to connected client
func (c *Client) Write() {

	defer func() {
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.MsgChan:
			if !ok {
				// manager closed the send channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			}

			// write message to connected client
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				c.Conn.Close()
			}
		}
	}
}
