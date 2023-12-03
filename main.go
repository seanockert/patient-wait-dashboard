package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Date  int64 `json:"date"`
	Level int   `json:"level"`
}

type PatientMessage struct {
	Patients []Message `json:"patients"`
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan PatientMessage)
)

func serveWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	clients[conn] = true

	// Listen for new messages from the client
	for {
		var patientMsg PatientMessage
		err := conn.ReadJSON(&patientMsg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}

		broadcast <- patientMsg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		patientMsg := <-broadcast

		// Send the message to every connected client
		for client := range clients {
			err := client.WriteJSON(patientMsg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {

	// Serve static assets
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Set up the HTTP route for both WebSocket and regular HTTP connections
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") == "websocket" {
			serveWebSocket(w, r)
		} else {
			http.ServeFile(w, r, "index.html")
		}
	})

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "admin.html")
	})

	// Handle incoming WebSocket messages
	go handleMessages()

	fmt.Println("Server is running on 0.0.0.0:8090")
	http.ListenAndServe("0.0.0.0:8090", nil)
}
