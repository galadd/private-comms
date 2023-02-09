package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	username := conn.LocalAddr().String()
	fmt.Println("Connected to server as:", username)

	messageChan := make(chan string)
	go func() {
		for {
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Connection closed")
				return
			}
			messageChan <- message
		}
	}()

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter message to send: ")
			message, _ := reader.ReadString('\n')
			if message == "\n" {
				continue
			}
			fmt.Fprint(conn, message)
		}
	}()

	for {
		select {
		case message := <-messageChan:
			fmt.Println()
			fmt.Println("Received from server:", message)
			fmt.Print("Enter message to send: ")
		}
	}
}
