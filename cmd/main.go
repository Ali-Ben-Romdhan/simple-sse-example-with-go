package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Client struct to store connected clients
type Client struct {
	ID      string
	Message chan string
}

var (
	mu          sync.Mutex
	clients     = make(map[*http.Request]Client)
	idCounter   = 0
	eventID     = 0
	eventIDFile = "eventID.txt"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/events", handleSSE)
	r.Post("/generate", storeEventID)
	r.Get("/generate", renderTemplate("generator.html"))
	r.Get("/display", renderTemplate("display.html"))

	http.Handle("/", r)

	http.ListenAndServe(":3000", r)
	fmt.Println("Server is running on :3000")
}

func handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set headers to indicate that this is an SSE endpoint
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create a channel to send messages to the client
	messageChan := make(chan string)

	// Register the client's message channel
	registerClient(r, messageChan)

	// Close the message channel when the client disconnects
	defer unregisterClient(r)

	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", message)
			w.(http.Flusher).Flush()
		}
	}
}

func registerClient(r *http.Request, messageChan chan string) {
	mu.Lock()
	defer mu.Unlock()

	idCounter++
	client := Client{
		ID:      fmt.Sprintf("client%d", idCounter),
		Message: messageChan,
	}

	clients[r] = client
	fmt.Printf("Client %s connected\n", client.ID)
}

func unregisterClient(r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	client, ok := clients[r]
	if ok {
		close(client.Message)
		delete(clients, r)
		fmt.Printf("Client %s disconnected\n", client.ID)
	}
}

func sendMessage(message string) {
	for _, client := range clients {
		client.Message <- message
	}
}

// triggerEvent simulates some event trigger in your application
func triggerEvent() {
	eventID, err := readEventIDFromFile()
	if err != nil {
		fmt.Println(err)
	}
	message := fmt.Sprintf("%d", eventID)
	sendMessage(message)
}

func storeEventID(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Increment eventID
	eventID++
	writeEventIDToFile(eventID)
	triggerEvent()
	w.Write([]byte(""))
}

func renderTemplate(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the absolute path to the directory containing the Go source file
		_, currentFile, _, _ := runtime.Caller(0)
		dir := filepath.Dir(currentFile)

		// Build the absolute path to the HTML file
		htmlFilePath := filepath.Join(dir, templateName)

		// Parse the HTML template
		tmpl, err := template.ParseFiles(htmlFilePath)
		if err != nil {
			http.Error(w, "Internal Server Error parsing template", http.StatusInternalServerError)
			return
		}

		// Execute the template, passing nil as data since we don't have any dynamic content
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func readEventIDFromFile() (int, error) {
	content, err := ioutil.ReadFile(eventIDFile)
	if err != nil {
		fmt.Println("Error reading eventID from file:", err)
		return 0, err
	}

	// Parse the content as an integer
	var parsedEventID int
	_, err = fmt.Sscanf(string(content), "%d", &parsedEventID)
	if err != nil {
		fmt.Println("Error parsing eventID from file:", err)
		return 0, err
	}
	return parsedEventID, nil
}

// writeEventIDToFile function to write eventID to a file
func writeEventIDToFile(eventID int) {
	err := ioutil.WriteFile(eventIDFile, []byte(fmt.Sprintf("%d", eventID)), 0644)
	if err != nil {
		fmt.Println("Error writing eventID to file:", err)
	}
}
