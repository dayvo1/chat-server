package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Chat Server — Step 1: WebSocket Basics
//
// HTTP is a request/response protocol — the client asks, the server answers, connection closes.
// WebSocket is different — it's a persistent two-way connection. Once established,
// either side can send messages at any time without the other side asking first.
// That's what makes real-time chat possible.
//
// How a WebSocket connection starts:
//   1. Client sends a normal HTTP request with a special "Upgrade" header
//   2. Server upgrades the connection from HTTP to WebSocket
//   3. The connection stays open — both sides can now send messages freely
//
// The upgrader:
//   var upgrader = websocket.Upgrader{
//       CheckOrigin: func(r *http.Request) bool { return true },
//   }
//   — gorilla/websocket handles the upgrade for you
//   — CheckOrigin controls which domains can connect. true = allow all (fine for development)
//
// Upgrading a connection:
//   conn, err := upgrader.Upgrade(w, r, nil)
//   — turns the HTTP connection into a WebSocket connection
//   — conn is the WebSocket connection you read/write on
//
// Reading a message:
//   messageType, message, err := conn.ReadMessage()
//   — blocks until a message arrives
//   — messageType is usually websocket.TextMessage
//   — message is []byte
//
// Writing a message:
//   err := conn.WriteMessage(websocket.TextMessage, []byte("hello"))
//
// Tasks:
//
// 1. Declare a package-level upgrader variable
//
// 2. Write a wsHandler(w http.ResponseWriter, r *http.Request) function:
//    — upgrade the connection
//    — if upgrade fails, log the error and return
//    — defer conn.Close()
//    — loop forever:
//        read a message from conn
//        if error, break (client disconnected)
//        log the message: log.Printf("received: %s", message)
//        echo it back: conn.WriteMessage(messageType, message)
//
// 3. In main:
//    — register wsHandler on "/ws"
//    — log "server running on :8080"
//    — ListenAndServe on :8080
//
// To test (you need a WebSocket client — use wscat):
//   npm install -g wscat
//   wscat -c ws://localhost:8080/ws
//   — type a message and hit enter, server should echo it back

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()

		if err != nil {
			return
		}

		log.Printf("received: %s", message)
		errWriting := conn.WriteMessage(messageType, message)

		if errWriting != nil {
			return
		}
	}
}
func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
