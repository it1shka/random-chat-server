package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: time.Second * 5,
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	connection, _ := upgrader.Upgrade(w, r, nil)
	processConnectedUser(connection)
}

func main() {
	setInfiniteLoop(time.Second, matchmaking)
	http.HandleFunc("/", websocketHandler)
	PORT := ":8080"
	fmt.Printf("Preparing for listening on port %s\n", PORT)
	http.ListenAndServe(PORT, nil)
}
