package main

import (
	"os"
	"fmt"

	"github.com/galadd/private-network/network"
	"gopkg.in/AlecAivazis/survey.v1"
)

func main() {
	// fmt.Println("Welcome to the Private Network")
	// fmt.Println("You can either")
	// fmt.Println("1. Create rsa key pairs ")
	// fmt.Println("Note: It is recommended to create rsa keys before sending messages")
	choice := ""
    prompt := &survey.Select{
        Message: "Choose an operation:",
        Options: []string{
			"Start private-comms", 
			"exit",
		},
    }
    survey.AskOne(prompt, &choice, nil)

	switch choice {
	case "Start private-comms":
		privateComms()
	case "exit":
		os.Exit(0)
	}
}

func privateComms() {
	init := ""
	prompt := &survey.Select{
		Message: "Are you the initiator or responder?",
		Options: []string{
			"initiator",
			"responder",
		},
	}
	survey.AskOne(prompt, &init, nil)
	switch init {
	case "initiator":
		network.ServerMain()
	case "responder":
		network.ClientMain()
	}
}