package chathub

import (
	"fmt"
	"log"
	"net/http"
)

// Handler struct
type Handler struct {
	Manager *Manager
}

// WsHandler handler
func (h *Handler) WsHandler(res http.ResponseWriter, req *http.Request) {
	sock, err := h.Manager.Upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Fatal(err)
	}

	id := req.Header.Get("Sec-Websocket-Key")
	fmt.Println(id)

	var client Client
	client.ID = id
	client.Conn = sock
	client.MsgChan = make(chan []byte)
	client.Room = make(map[string]bool)
	client.Manager = h.Manager

	h.Manager.Register <- &client

	// Consume message
	go client.Consume()

	// Publish message
	go client.Publish()
}
