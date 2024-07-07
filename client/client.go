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

	u := fmt.Sprintf("ws://%s", uri)
	log.Printf("connecting to %s", u)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

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

	// Send announce message
	err = wc.announcePresence()
	if err != nil {
		log.Println("announce:", err)
	}

	// Periodically send ping to keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then waiting (with timeout) for server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}

// discover all the clients in the network which have called the -connect() function
func Discover() {
	/*
		the connect() function does the job of discovering all the clients
		this function will fetch the connected clients and display
	*/

}

// function to send the file to the desired client
func Send(receiver string, inputFilePath string) {
	//function will require unique identifier of a client to send it
}

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
