package network

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"crypto/rsa"

	auth "github.com/galadd/private-comms/authentication"
	e "github.com/galadd/private-comms/encryption"
)

func ServerMain(myPrivateKey *rsa.PrivateKey, respPublicKey *rsa.PublicKey) {
	fmt.Println("Starting server...")
	ln, err := net.Listen("tcp", ":4357")
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
		go handleConnection(conn, myPrivateKey, respPublicKey)
	}
}

func handleConnection(conn net.Conn, myPrivateKey *rsa.PrivateKey, respPublicKey *rsa.PublicKey) {
	keyPair, err := e.GenerateKeypair()
	if err != nil {
		fmt.Println(err)
		return
	}

	username := conn.RemoteAddr().String()
	fmt.Println("Accepted connection from:", username)	
	defer conn.Close()

	conn.Write(keyPair.PublicKey[:])
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
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
	key, err := e.DH(keyPair.PrivateKey[:], buf)
	if err != nil {
		fmt.Println(err)
		return
	}

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
			signature := buf[:256]
			buf = buf[256:]

			err = auth.VerifySignature(respPublicKey, buf, signature)
			if err != nil {
				fmt.Println(err)
				return
			}

			decrypt, _ := e.Decrypt(buf, key)

			messageChan <- decrypt
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			fmt.Print("=> ")	
			message, _ := reader.ReadString('\n')
			fmt.Print("\033[1A\033[2K")
			fmt.Print("You: " + message)

			byteMessage := []byte(message)
			ciphertext, _ := e.Encrypt(byteMessage, key)


			signature, err := auth.SignMessage(myPrivateKey, ciphertext)
			if err != nil {
				fmt.Println(err)
				return
			}

			ciphertext = append(signature, ciphertext...)
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
