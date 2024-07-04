package client

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn *websocket.Conn
}

func (wc *WebSocketClient) connect(uri string) {
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

func (wc *WebSocketClient) announcePresence() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}

	var ip string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	if ip == "" {
		return fmt.Errorf("unable to determine local IP address")
	}

	message := fmt.Sprintf(`{"type":"announce","hostname":"%s","ip":"%s"}`, hostname, ip)
	err = wc.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}

	log.Printf("sent: %s", message)
	return nil
}

func main() {
	// Replace with your WebSocket server URI
	serverURI := "localhost:8080/ws"
	wc := WebSocketClient{}

	// Connect to WebSocket server
	wc.connect(serverURI)
}
