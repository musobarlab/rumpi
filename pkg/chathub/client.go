package chathub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/musobarlab/rumpi/pkg/jwt"
)

const (

	// WriteWait Time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// PongWait Time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// PongPeriod should less than PongWait
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 1024
)

// Client model
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	MsgChan  chan []byte
	Manager  *Manager
	Room     map[string]bool
	IsOnline bool
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

// Consume function will read incoming message from client
func (c *Client) Consume() {

	// if function finish, remove client
	// and also close connection
	defer func() {
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

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

		if message.MessageType == AuthMessage {
			if message.AuthKey != c.Manager.AuthKey {
				// auth is invalid, then remove the client that match with its ID
				c.Conn.Close()
				c.Manager.DeleteClient(c)
			} else {
				jwtClaimResult := c.Manager.JwtService.Validate(context.Background(), message.Token)
				if jwtClaimResult.Error != nil {
					c.Conn.Close()
					c.Manager.DeleteClient(c)
				} else {
					jwtClaim := jwtClaimResult.Data.(*jwt.Claim)
					fmt.Println("-------------------------")
					fmt.Println(jwtClaim.User.Email)
					// auth success, send information to AuthSuccess's Manager
					c.IsOnline = true
					if message.Username != "" {
						c.Username = message.Username
					}
					c.Manager.AuthSuccess <- c
				}
			}
		}

		message.From = c.Username
		message.Date = time.Now()

		// send message to Manager's IncomingMessage
		c.Manager.IncomingMessage <- &message

	}

}

// Publish function will write message to connected client
func (c *Client) Publish() {

	ticker := time.NewTicker(PingPeriod)

	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.MsgChan:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// manager closed the send channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// write message to connected client
			// err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			// if err != nil {
			// 	c.Conn.Close()
			// }

			writer, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			writer.Write(msg)

			// Add queued chat messages to the current websocket message.
			n := len(c.MsgChan)
			for i := 0; i < n; i++ {
				writer.Write([]byte{'\n'})
				writer.Write(<-c.MsgChan)
			}

			err = writer.Close()
			if err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
