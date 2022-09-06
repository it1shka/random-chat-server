package main

import "github.com/gorilla/websocket"

var queue = NewQueue()
var dispatcher = NewDict[*websocket.Conn, *websocket.Conn]()

func matchmaking() {
	queue.DoMatchmaking(func(a, b *websocket.Conn) {
		dispatcher.Set(a, b)
		dispatcher.Set(b, a)
		MessageFoundPartner(a)
		MessageFoundPartner(b)
	}, 100)
}

func processConnectedUser(conn *websocket.Conn) {
	queue.Append(conn)
	defer func() {
		paired := dispatcher.Get(conn)
		dispatcher.Delete(conn)
		dispatcher.Delete(paired)
		conn.Close()
		if paired != nil {
			MessageEnd(paired)
		}
	}()
	handlingConnection(conn)
}

func handlingConnection(conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil || mt == websocket.CloseMessage {
			break
		}
		sendTo := dispatcher.Get(conn)
		if sendTo != nil {
			MessageText(sendTo, string(message))
		}
	}
}
