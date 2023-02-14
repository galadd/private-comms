package main

import (
	"fmt"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"

	"github.com/galadd/private-network/network"
	"gopkg.in/AlecAivazis/survey.v1"
	"github.com/joho/godotenv"

	auth "github.com/galadd/private-network/authentication"
)

var (
	myAuthPrivateKey  *rsa.PrivateKey
	respAuthPublicKey *rsa.PublicKey

	hexEncodedPriv    string
	hexEncodedPub     string

	err 			  error
)

func main() {
	fmt.Println("Start the program if you have a keypair and responder's public key in .env file")
	usage := ""
	prompt := &survey.Select{
		Message: "Select",
		Options: []string{
			"Generate a new keypair",
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
		respHexPub := os.Getenv("RESPONDER_PUB_KEY")
		respAuthPublicKey, err = decodeHexPub(respHexPub)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		myHexPriv := os.Getenv("HEX_ENCODED_PRIVATE")
		myAuthPrivateKey, err = decodeHexPriv(myHexPriv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		start := ""
		prompt = &survey.Select{
			Message: "",
			Options: []string{
				"Start",
				"Exit",
			},
		}
		survey.AskOne(prompt, &start, nil)

		switch start {
		case "Start":
			operation()
		case "Exit":
			os.Exit(0)
		}

	case "Use an existing keypair":
		err = godotenv.Load(".env")
		if err != nil {
			fmt.Println(err)
			return
		}
		myHexPriv := os.Getenv("HEX_ENCODED_PRIVATE")
		myAuthPrivateKey, err = decodeHexPriv(myHexPriv)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		respHexPub := os.Getenv("RESPONDER_PUB_KEY")
		respAuthPublicKey, err = decodeHexPub(respHexPub)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		operation()
		

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
		network.ServerMain(myAuthPrivateKey, respAuthPublicKey)
	case "responder":
		network.ClientMain(myAuthPrivateKey, respAuthPublicKey)
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

func decodeHexPriv(hexEncodedPriv string) (*rsa.PrivateKey, error) {
	// convert hex encoded string to byte slice
	decoded, err := hex.DecodeString(hexEncodedPriv)
	if err != nil {
		return nil, err
	}

	// load private key
	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func decodeHexPub(hexEncodedPub string) (*rsa.PublicKey, error) {
	// convert hex encoded string to byte slice
	decoded, err := hex.DecodeString(hexEncodedPub)
	if err != nil {
		return nil, err
	}

	// load public key
	block, _ := pem.Decode(decoded)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

func write(hexEncoded, hexEncodedPub string) {
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	envVars := map[string]string{
		"HEX_ENCODED_PRIVATE": hexEncoded,
		"HEX_ENCODED_PUBLIC":   hexEncodedPub,
	}

	for key, value := range envVars {
		_, err = file.WriteString(key + "=" + value + "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}