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

func ClientMain(myPrivateKey *rsa.PrivateKey, respPublicKey *rsa.PublicKey, ip string) {
	conn, err := net.Dial("tcp", ip + ":4357")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return
	}
	defer conn.Close()

	keyPair, err := e.GenerateKeypair()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	conn.Write(keyPair.PublicKey[:])

	username := conn.LocalAddr().String()
	fmt.Println("Connected to server as", username)

	// receive public key from server
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

			// verify signature
			err  = auth.VerifySignature(respPublicKey, buf, signature)
			if err != nil {
				fmt.Println(err)
				return
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
