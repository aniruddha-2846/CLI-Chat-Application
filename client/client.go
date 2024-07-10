package client

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// how will you identify every unique client?
type clientInfo struct {
	userName  string
	ipAddress string
}

type WebSocketClient struct {
	conn *websocket.Conn
}

// this function connects the client to the server
func (wc *WebSocketClient) Connect(uri string) {

	/*
		send an http request to the websocket server which will then upgrade it to
		a websocket connection.

	*/
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// u := fmt.Sprintf("ws://%s", uri)
	u := fmt.Sprintf("ws://%s/ws", uri)
	log.Printf("connecting to %s", u)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	wc.conn = conn
	log.Println("connected to", u)

	done := make(chan struct{})

	// Handle incoming messages in a separate goroutine
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case <-interrupt:
	// 		log.Println("interrupt")
	// 		err := wc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 		if err != nil {
	// 			log.Println("write close:", err)
	// 			return
	// 		}
	// 		select {
	// 		case <-done:
	// 		case <-time.After(time.Second):
	// 		}
	// 		wc.conn.Close()
	// 		return
	// 	}
	// }
	// Send announce message
	// err = wc.announcePresence()
	// if err != nil {
	// 	log.Println("announce:", err)
	// }

}

// discover all the clients in the network which have called the -connect() function
func Discover() {
	/*
		the connect() function does the job of discovering all the clients
		this function will fetch the connected clients and display
	*/

}

// function to send the file to the desired client
// func Send(sender *WebSocketClient, receiver *WebSocketClient, inputFilePath string) {
// 	//function will require unique identifier of a client to send it
// 	//1. open the file and read it to a buffer
// 	//2. send the file
// }

// cuts off the connection of this client with the websocket server
func (client *WebSocketClient) Disconnect() {
	deadline := time.Now().Add(time.Minute)
	err := client.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		deadline,
	)
	if err != nil {
		fmt.Printf("Error in WriteControl: %s", err)
		return
	}

	// Set deadline for reading the next message
	err = client.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		fmt.Printf("Error in SetReadDeadline: %s", err)
		return
	}
	// Read messages until the close message is confirmed
	for {
		_, _, err = client.conn.NextReader()
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			break
		}
		if err != nil {
			break
		}
	}
	// Close the TCP connection
	err = client.conn.Close()
	if err != nil {
		fmt.Printf("Error closing the TCP connection: %s", err)
		return
	}
}

// !  new new new new nwe new
// package client

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"time"

// 	"github.com/gorilla/websocket"
// )

// type WebSocketClient struct {
// 	conn     *websocket.Conn
// 	url      string
// 	mutex    sync.Mutex
// 	isClosed bool
// }

// func NewWebSocketClient(url string) *WebSocketClient {
// 	return &WebSocketClient{
// 		url: url,
// 	}
// }

// func (wc *WebSocketClient) Connect() {
// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt)

// 	u := fmt.Sprintf("ws://%s/ws", wc.url)
// 	log.Printf("connecting to %s", u)

// 	dialer := websocket.DefaultDialer
// 	conn, _, err := dialer.Dial(u, nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 		return
// 	}
// 	wc.mutex.Lock()
// 	wc.conn = conn
// 	wc.isClosed = false
// 	wc.mutex.Unlock()

// 	log.Println("connected to", u)

// 	done := make(chan struct{})

// 	go func() {
// 		defer close(done)
// 		for {
// 			_, message, err := conn.ReadMessage()
// 			if err != nil {
// 				log.Println("read:", err)
// 				wc.Disconnect()
// 				return
// 			}
// 			log.Printf("recv: %s", message)
// 		}
// 	}()

// 	for {
// 		select {
// 		case <-done:
// 			return
// 		case <-interrupt:
// 			log.Println("interrupt")
// 			wc.Disconnect()
// 			return
// 		}
// 	}
// }

// func (wc *WebSocketClient) Disconnect() {
// 	wc.mutex.Lock()
// 	defer wc.mutex.Unlock()

// 	if wc.isClosed {
// 		return
// 	}

// 	deadline := time.Now().Add(time.Minute)
// 	err := wc.conn.WriteControl(
// 		websocket.CloseMessage,
// 		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
// 		deadline,
// 	)
// 	if err != nil {
// 		fmt.Printf("Error in WriteControl: %s", err)
// 		return
// 	}

// 	err = wc.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
// 	if err != nil {
// 		fmt.Printf("Error in SetReadDeadline: %s", err)
// 		return
// 	}
// 	for {
// 		_, _, err = wc.conn.NextReader()
// 		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
// 			break
// 		}
// 		if err != nil {
// 			break
// 		}
// 	}
// 	err = wc.conn.Close()
// 	if err != nil {
// 		fmt.Printf("Error closing the TCP connection: %s", err)
// 		return
// 	}

// 	wc.isClosed = true
// 	log.Println("disconnected from server")
// }

// func (wc *WebSocketClient) Discover() {
// 	wc.mutex.Lock()
// 	defer wc.mutex.Unlock()

// 	if wc.conn == nil {
// 		log.Println("Not connected")
// 		return
// 	}

// 	err := wc.conn.WriteJSON(map[string]interface{}{
// 		"command": "discover",
// 	})
// 	if err != nil {
// 		log.Println("Error sending discover command:", err)
// 	}
// }

// func (wc *WebSocketClient) Send(receiver, inputFilePath string) {
// 	wc.mutex.Lock()
// 	defer wc.mutex.Unlock()

// 	if wc.conn == nil {
// 		log.Println("Not connected")
// 		return
// 	}

// 	err := wc.conn.WriteJSON(map[string]interface{}{
// 		"command":  "send",
// 		"receiver": receiver,
// 		"path":     inputFilePath,
// 	})
// 	if err != nil {
// 		log.Println("Error sending file:", err)
// 	}
// }
