package main

import (
	"fmt"
	"net"
	"bufio"
	"os"

	e "github.com/galadd/private-network/encryption"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4357")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()

	username := conn.LocalAddr().String()
	fmt.Println("Connected to server as", username)

	key := []byte("01234567890123456789012345678901")

	messageChan := make(chan []byte)
	go func() {
		for {
			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}

			for i := len(buf) - 1; i >= 0; i-- {
				if buf[i] != 0 {
					break
				}
				buf = buf[:i]
			}

			decrypt, _ := e.Decrypt(buf, key)

			messageChan <- decrypt
		}
	}()

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("=> ")
			message, _ := reader.ReadString('\n')
			fmt.Print("\033[1A\033[2K")
			fmt.Print("You: " + message)

			byteMessage := []byte(message)
			ciphertext, _ := e.Encrypt(byteMessage, key)

			conn.Write(ciphertext)
		}
	}()

	for {
		select {
		case message := <-messageChan:
			fmt.Printf("%v: %v", username, string(message))
			fmt.Print("=> ")
		}
	}
}
