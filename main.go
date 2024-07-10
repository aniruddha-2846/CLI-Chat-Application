package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"webfileshare/client"
)

// fileshare -connect()
// fileshare -discover()
// fileshare -send [receiver_name] [input_file_path]
// fileshare -disconnect()

func main() {
	url := "localhost:8080"
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	connectFlag := flag.Bool("connect", false, "Connect to the server")
	flag.Parse()

	if *connectFlag {
		wsClient := &client.WebSocketClient{}
		wsClient.Connect(url)
		fmt.Println("Connected to server.")
	} else {
		flag.Usage()
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		switch input {
		case "-discover":
			client.Discover()
		case "-send":
			fmt.Println("Usage: -send receiver_name input_file_path")
		case "-disconnect":
			wsClient := &client.WebSocketClient{}
			wsClient.Disconnect()
			fmt.Println("Disconnected from server.")
			return
		default:
			fmt.Println("Invalid command. Available commands: -discover, -send, -disconnect")
		}
	}
}

// func main() {
// 	url := "localhost:8080"

// 	// Define flags
// 	connectFlag := flag.Bool("connect", false, "Connect to the server")
// 	disconnectFlag := flag.Bool("disconnect", false, "Disconnect from the server")
// 	discoverFlag := flag.Bool("discover", false, "Discover all clients in your network")
// 	sendFlag := flag.String("send", "", "Send file to receiver (usage: -send receiver_name input_file_path)")

// 	// Parse flags
// 	flag.Parse()

// 	wsClient := client.NewWebSocketClient(url)
// 	connectChan := make(chan bool)
// 	disconnectChan := make(chan bool)
// 	discoverChan := make(chan bool)
// 	sendChan := make(chan [2]string)

// 	go func() {
// 		for {
// 			select {
// 			case <-connectChan:
// 				wsClient.Connect()
// 			case <-disconnectChan:
// 				wsClient.Disconnect()
// 			case <-discoverChan:
// 				wsClient.Discover()
// 				// case args := <-sendChan:
// 				// 	wsClient.Send()
// 			}
// 		}
// 	}()

// 	fmt.Println("Welcome to my CLI Filesharing application")

// 	if *connectFlag {
// 		connectChan <- true
// 	} else if *disconnectFlag {
// 		disconnectChan <- true
// 	} else if *discoverFlag {
// 		discoverChan <- true
// 	} else if *sendFlag != "" {
// 		args := flag.Args()
// 		if len(args) < 2 {
// 			fmt.Println("Usage: -send receiver_name input_file_path")
// 		} else {
// 			sendChan <- [2]string{args[0], args[1]}
// 		}
// 	} else {
// 		log.Println("Please enter a valid command!")
// 		flag.Usage()
// 	}
// }
