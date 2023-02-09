package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
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
	fmt.Println("New client connected:", username)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", username)
			return
		}
		if strings.TrimSpace(message) == "" {
			continue
		}
		fmt.Println(username + ": " + message)
	}
}
