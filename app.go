package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Date  int64 `json:"date"`
	Level int   `json:"level"`
}

type PatientMessage struct {
	Patients []Message `json:"patients"`
}

type cacheControlFile struct {
	http.File
}

// Cache-Control max age in seconds
const cacheMaxAge = 300 // 31536000

func (f cacheControlFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f cacheControlFile) Stat() (os.FileInfo, error) {
	info, err := f.File.Stat()
	if err != nil {
		return nil, err
	}

	return fileInfoWithCacheHeaders{info}, nil
}

type fileInfoWithCacheHeaders struct {
	os.FileInfo
}

func (f fileInfoWithCacheHeaders) ModTime() time.Time {
	return time.Unix(0, 0) // Always use the Unix epoch time for simplicity
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	host := "0.0.0.0:" + port

	fmt.Println("Server is running on:", host)
	panic(http.ListenAndServe(host, nil))
}

func serveStaticFileWithCache(w http.ResponseWriter, r *http.Request, filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error opening file: %v", err)
		return
	}
	defer file.Close()

	// Get file information
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error getting file information: %v", err)
		return
	}

	// Set caching headers
	modTime := fileInfo.ModTime()
	etag := fmt.Sprintf("%x", modTime.UnixNano())
	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=31536000") // One year cache duration
	w.Header().Set("Expires", modTime.AddDate(1, 0, 0).UTC().Format(http.TimeFormat))

	// Check If-None-Match header for caching
	if match := r.Header.Get("If-None-Match"); match == etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Serve the content
	http.ServeContent(w, r, fileInfo.Name(), modTime, file)
}
