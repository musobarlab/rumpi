package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

type ClientModel struct{
	Username string `json:"username"`
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

var registeredClients = make(map[string]*ClientModel)

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var manager = ClientManager{
	clients:    make(map[*Client]bool),
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
			jsonMessage, _ := json.Marshal(&Message{Sender: conn.id, Content: "has join the chat"})
			manager.send(jsonMessage, conn)
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Sender: conn.id, Content: "has leave the chat"})
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
		manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprintf(res, "hell yeah!")
}

func registerHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	var clientModel ClientModel
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&clientModel)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := registeredClients[clientModel.Username]; ok {
		http.Error(res, "Username already exist", http.StatusBadRequest)
		return
	}

	registeredClients[clientModel.Username] = &clientModel

	clientModelRes, err := json.Marshal(clientModel)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(res, string(clientModelRes))
}

func websocketHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	id := req.Header.Get("Sec-Websocket-Key")
	fmt.Println(id)

	client := &Client{id: id, socket: conn, send: make(chan []byte)}

	manager.register <- client

	go client.read()
	go client.write()
}

func main() {
	fmt.Println("Starting application...")
	go manager.start()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/join", registerHandler).Methods("POST")
	router.HandleFunc("/ws", websocketHandler)

	log.Fatal(http.ListenAndServe(":9000", router))
}
