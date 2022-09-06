package main

import (
	"github.com/gorilla/websocket"
)

// communication channel system
var channels = NewDict[*websocket.Conn, *websocket.Conn]()

func CreateChannel(a, b *websocket.Conn) {
	channels.Set(a, b)
	channels.Set(b, a)
	MessageFoundPartner(a)
	MessageFoundPartner(b)
}

func DeleteChannel(conn *websocket.Conn) {
	paired := channels.Get(conn)
	channels.Delete(conn)
	channels.Delete(paired)
	if paired != nil {
		MessageEnd(paired)
	}
}

// matchmaking logic (main app logic)
var queue = NewQueue()

func matchmaking() {
	queue.DoMatchmaking(CreateChannel, 100)
}

func processConnectedUser(conn *websocket.Conn) {
	defer func() {
		queue.Delete(conn)
		DeleteChannel(conn)
		conn.Close()
	}()

	handlingConnection(conn)
}

func handlingConnection(conn *websocket.Conn) {
	for {
		var data Json
		if err := conn.ReadJSON(data); err != nil {
			break
		}

		switch data["type"] {
		// find next partner
		case "next":
			DeleteChannel(conn)
			queue.Append(conn)

		// send message to partner
		case "message":
			if partner := channels.Get(conn); partner != nil {
				message := data["message"].(string)
				MessageText(partner, message)
			}

		}
	}
}
