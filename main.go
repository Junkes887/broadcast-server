package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

var (
	clients   = make(map[net.Conn]bool)
	broadcast = make(chan string)
)

func main() {
	operation := flag.String("operation", "", "start, connect")
	port := flag.Int("port", 8080, "server port")
	username := flag.String("username", "You", "your username")
	flag.Parse()

	switch *operation {
	case "start":
		start(*port)
	case "connect":
		connect(*port, *username)
	default:
		fmt.Println("Text for --help to list operations")
	}
}

func start(port int) {
	fmt.Printf("Starting server on port %d...", port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer ln.Close()

	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				_, err := fmt.Fprint(client, msg)
				if err != nil {
					fmt.Println("Error broadcasting message:", err)
					delete(clients, client)
					client.Close()
				}
			}
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	clients[conn] = true

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			delete(clients, conn)
			return
		}
		broadcast <- message
	}
}

func connect(port int, username string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {

				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Print(message)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		message, _ := reader.ReadString('\n')
		_, err := fmt.Fprint(conn, fmt.Sprintf("%s: %s", username, message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}
