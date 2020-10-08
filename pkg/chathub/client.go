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

// Client model
type Client struct {
	*websocket.Conn
	ID       string
	Username string
	MsgChan  chan []byte
	Manager  *Manager
	Room     map[string]bool
	IsOnline bool
	sync.RWMutex
}

//joinRoom function will push new room to the map rooms
func (client *Client) joinRoom(name string) {
	client.Lock()
	client.Room[name] = true
	client.Unlock()
}

//leaveRoom function will delete room by specific key from map rooms
func (client *Client) leaveRoom(name string) {
	client.Lock()
	delete(client.Room, name)
	client.Unlock()
}

// Consume function will read incoming message from client
func (c *Client) Consume() {

	// if function finish, remove client
	// and also close connection
	defer func() {
		c.Manager.ExitedClient <- c
		c.Close()
	}()

	c.SetReadLimit(MaxMessageSize)
	c.SetReadDeadline(time.Now().Add(PongWait))
	c.SetPongHandler(func(string) error {
		c.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		// read message from client
		_, msg, err := c.ReadMessage()
		if err != nil {

			c.Manager.ExitedClient <- c
			c.Close()
			break
		}

		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			c.Manager.ExitedClient <- c
			c.Close()
			break
		}

		if message.MessageType == AuthMessage {
			if message.AuthKey != c.Manager.AuthKey {
				// auth is invalid, then remove the client that match with its ID
				c.Close()
				c.Manager.deleteClient(c.Username)
				break
			}

			jwtClaimResult := c.Manager.JwtService.Validate(context.Background(), message.Token)
			if jwtClaimResult.Error != nil {
				// this error scope, represent error from token expired or token format error
				// if error, send message with 'AuthFail' type, to tell client to close its connection
				// and remove its cookie or localstorage
				message := &Message{MessageType: AuthFail, Content: jwtClaimResult.Error.Error(), Date: time.Now()}
				msg, _ := json.Marshal(message)

				c.MsgChan <- msg

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
		c.Close()
	}()

	for {
		select {
		case msg, ok := <-c.MsgChan:
			c.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// manager closed the send channel.
				c.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// write message to connected client
			// err := c.WriteMessage(websocket.TextMessage, msg)
			// if err != nil {
			// 	c.Close()
			// }

			writer, err := c.NextWriter(websocket.TextMessage)
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
			c.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
