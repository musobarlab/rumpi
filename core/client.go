package core

import (
	"encoding/json"

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

		message.Sender = c.ID

		// send broadcast
		c.Manager.BroadCast <- &message

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
