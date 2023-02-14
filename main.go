package main

import (
	"fmt"
	// "crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"

	"github.com/galadd/private-network/network"
	"gopkg.in/AlecAivazis/survey.v1"

	auth "github.com/galadd/private-network/authentication"
)

var (
		hexEncodedPriv  string
		hexEncodedPub   string

		err 			error
)

func main() {
	usage := ""
	prompt := &survey.Select{
		Message: "Select",
		Options: []string{
			"Create a new keypair",
			"Use an existing keypair",
			"Exit",
		},
	}
	survey.AskOne(prompt, &usage, nil)

	switch usage {
	case "Generate a new keypair":
		hexEncodedPriv, hexEncodedPub, err = generate()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		write(hexEncodedPriv, hexEncodedPub)

	case "Use an existing keypair":
		

	case "Exit":
		os.Exit(0)
	}
}

func operation() {
	choice := ""
	prompt := &survey.Select{
		Message: "Are you the initiator or responder?",
		Options: []string{
			"initiator",
			"responder",
			"exit",
		},
	}
	survey.AskOne(prompt, &choice, nil)
	switch choice {
	case "initiator":
		network.ServerMain()
	case "responder":
		network.ClientMain()
	case "exit":
		os.Exit(0)
	}
}

func generate() (string, string, error) {
	privateKey, publicKey, err := auth.GenerateKeyPair()
	if err != nil {
		return "", "", err
	}

	// encode private key to pem format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type: "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	pemEncodedPriv := pem.EncodeToMemory(privateBlock)

	// encode public key to pem format
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	publicBlock := &pem.Block{
		Type: "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	pemEncodedPub := pem.EncodeToMemory(publicBlock)

	// convert to hex encoded string
	hexEncodedPriv := hex.EncodeToString(pemEncodedPriv)
	hexEncodedPub := hex.EncodeToString(pemEncodedPub)

	return hexEncodedPriv, hexEncodedPub, nil
}

func write(hexEncoded, hexEncodedPub string) {
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	envVars := map[string]string{
		"hexEncoded-Private": hexEncoded,
		"hexEncoded-Public":   hexEncodedPub,
	}

	for key, value := range envVars {
		_, err = file.WriteString(key + "=" + value + "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}