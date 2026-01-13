package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Simple echo or status update logic would go here.
		// For this task, we can just keep the connection open or implement a broadcast later.
		// The user asked for HTMX+Websocket, usually HTMX uses WS to swap content.

		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("Received: %s", message)
	}
}
