package main

import (
	"os"
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting server...")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	username := conn.RemoteAddr().String()
	fmt.Println("Accepted connection from:", username)

	messageChan := make(chan string)
	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Closed connection from:", username)
				break
			}
			messageChan <- message
		}
		conn.Close()
	}()

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter message to send: ")
			message, _ := reader.ReadString('\n')
			if message == "\n" {
				continue
			}
			sendMessage := username + ": " + message
			fmt.Println(sendMessage)
			fmt.Fprint(conn, sendMessage)
		}
	}()

	for {
		select {
		case message := <-messageChan:
			fmt.Println()
			fmt.Println("Received from client:", message)
			fmt.Print("Enter message to send: ")
		}
	}
}

