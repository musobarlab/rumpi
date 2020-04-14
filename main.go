package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/musobarlab/go-websocket-chat/core"
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprintf(res, "hell yeah!")
}

func registerHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")

	var clientModel core.ClientModel
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&clientModel)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// if _, ok := registeredClients[clientModel.Username]; ok {
	// 	http.Error(res, "Username already exist", http.StatusBadRequest)
	// 	return
	// }

	// registeredClients[clientModel.Username] = &clientModel

	clientModelRes, err := json.Marshal(clientModel)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(res, string(clientModelRes))
}

func main() {
	fmt.Println("starting application...")
	manager := core.NewManager()

	wsHandler := core.Handler{Manager: manager}

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/join", registerHandler).Methods("POST")
	router.HandleFunc("/ws", wsHandler.WsHandler)

	// start client manager
	go manager.Run()

	fmt.Println("server running on port 9000")
	log.Fatal(http.ListenAndServe(":9000", router))
}
