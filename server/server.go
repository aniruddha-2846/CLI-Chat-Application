package server

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]string) // connected clients
var clientsLock sync.Mutex
var connectionCounter int

// Message object
type Message struct {
	Type     string `json:"type"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

func deleteClient(ws *websocket.Conn) {
	clientsLock.Lock()
	delete(clients, ws)
	connectionCounter--
	clientsLock.Unlock()
}

func broadcastMessage(sender *websocket.Conn, network string, msg map[string]interface{}) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	for client, clientNetwork := range clients {
		if client != sender && clientNetwork == network {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func handleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// Get the client's IP address and port number
	clientAddr := ws.RemoteAddr().(*net.TCPAddr)
	clientIP := clientAddr.IP.String()
	clientPort := clientAddr.Port

	connectionCounter++
	log.Printf("New connection from %s:%d, Total connections: %d", clientIP, clientPort, connectionCounter)

	network := r.URL.Query().Get("network")
	if network == "" {
		network = "default"
	}

	clientsLock.Lock()
	clients[ws] = network
	clientsLock.Unlock()

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			deleteClient(ws)
			break
		}
		log.Printf("Received message from %s: %v", network, msg)
		broadcastMessage(ws, network, msg)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	log.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
