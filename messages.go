package main

import "github.com/gorilla/websocket"

func MessageFoundPartner(conn *websocket.Conn) {
	conn.WriteJSON(Json{
		"type": "start",
	})
}

func MessageText(conn *websocket.Conn, text string) {
	conn.WriteJSON(Json{
		"type": "message",
		"text": text,
	})
}

func MessageEnd(conn *websocket.Conn) {
	conn.WriteJSON(Json{
		"type": "end",
	})
}
