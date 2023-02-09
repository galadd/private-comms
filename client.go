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

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		if message == "\n" {
			continue
		}
		fmt.Fprint(conn, message)
	}
}
