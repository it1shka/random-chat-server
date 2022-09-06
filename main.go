package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "This is a websocket application")
		return
	}
	processConnectedUser(connection)
}

func main() {
	setInfiniteLoop(time.Second, matchmaking)
	http.HandleFunc("/", websocketHandler)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	fmt.Printf("Preparing for listening on port %s\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
